package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
)

// EnrichCharacterInput contains input for enrichment
type EnrichCharacterInput struct {
	CharacterName string
	Inplace       bool // Update original file vs create enriched version
	Force         bool // Force refetch even if cached
}

// EnrichCharacterOutput contains enrichment results
type EnrichCharacterOutput struct {
	OutputPath string
	Fetched    map[string]int
	Cached     map[string]int
	Logs       []string
}

// EnrichCharacterUseCase handles fetching external API data for a character
type EnrichCharacterUseCase struct {
	characterRepo  repositories.CharacterRepository
	enrichmentRepo repositories.EnrichmentRepository
}

// NewEnrichCharacterUseCase creates a new EnrichCharacterUseCase
func NewEnrichCharacterUseCase(
	characterRepo repositories.CharacterRepository,
	enrichmentRepo repositories.EnrichmentRepository,
) *EnrichCharacterUseCase {
	return &EnrichCharacterUseCase{
		characterRepo:  characterRepo,
		enrichmentRepo: enrichmentRepo,
	}
}

// Execute enriches a character with external API data
func (uc *EnrichCharacterUseCase) Execute(ctx context.Context, input EnrichCharacterInput) (*EnrichCharacterOutput, error) {
	logs := []string{}
	start := time.Now()

	char, err := uc.characterRepo.FindByName(input.CharacterName)
	if err != nil {
		return nil, fmt.Errorf("character not found: %w", err)
	}

	logs = append(logs, fmt.Sprintf("Loaded character: %s", input.CharacterName))

	// Collect spell and equipment names
	spellNames := uc.collectSpellNames(char)
	equipNames := uc.collectEquipmentNames(char)
	logs = append(logs, fmt.Sprintf("Found %d spells, %d equipment items", len(spellNames), len(equipNames)))

	// Check if enriched file already exists
	fetched := map[string]int{"spells": 0, "equipment": 0}
	cached := map[string]int{"spells": 0, "equipment": 0}
	
	if !input.Inplace {
		outPath := filepath.Join("data", "enrichments", input.CharacterName+".json")
		if _, err := os.Stat(outPath); err == nil && !input.Force {
			logs = append(logs, "Enriched file already exists; skipping (use --force to override)")
			return &EnrichCharacterOutput{
				OutputPath: outPath,
				Fetched:    fetched,
				Cached:     cached,
				Logs:       logs,
			}, nil
		}
	}

	// Load caches
	spellCache, _ := uc.enrichmentRepo.LoadSpellCache()
	equipCache, _ := uc.enrichmentRepo.LoadEquipmentCache()

	// Prepare jobs for concurrent fetching
	var jobs []fetchJob

	// Create spell jobs
	for _, spellName := range spellNames {
		spellIndex := uc.sanitizeIndex(spellName)
		if _, ok := spellCache[spellIndex]; ok && !input.Force {
			cached["spells"]++
			continue
		}
		jobs = append(jobs, fetchJob{"spell", spellName, spellIndex})
	}

	// Create equipment jobs
	for _, equipName := range equipNames {
		equipIndex := uc.sanitizeIndex(equipName)
		if _, ok := equipCache[equipIndex]; ok && !input.Force {
			cached["equipment"]++
			continue
		}
		jobs = append(jobs, fetchJob{"equipment", equipName, equipIndex})
	}

	// Fetch concurrently with rate limiting
	results := uc.fetchConcurrently(ctx, jobs, logs)

	// Process results
	enrichedSpells := make(map[string]repositories.SpellInfo)
	enrichedEquipment := make(map[string]repositories.EquipmentInfo)

	for _, result := range results {
		if result.err != nil {
			logs = append(logs, fmt.Sprintf("Warning: failed to fetch %s %s: %v", result.kind, result.name, result.err))
			continue
		}

		if result.kind == "spell" && result.spell != nil {
			enrichedSpells[result.name] = *result.spell
			spellCache[result.index] = *result.spell
			fetched["spells"]++
		} else if result.kind == "equipment" && result.equip != nil {
			enrichedEquipment[result.name] = *result.equip
			equipCache[result.index] = *result.equip
			fetched["equipment"]++
		}
	}

	// Add cached items to output
	for _, spellName := range spellNames {
		spellIndex := uc.sanitizeIndex(spellName)
		if info, ok := spellCache[spellIndex]; ok {
			if _, exists := enrichedSpells[spellName]; !exists {
				enrichedSpells[spellName] = info
			}
		}
	}
	for _, equipName := range equipNames {
		equipIndex := uc.sanitizeIndex(equipName)
		if info, ok := equipCache[equipIndex]; ok {
			if _, exists := enrichedEquipment[equipName]; !exists {
				enrichedEquipment[equipName] = info
			}
		}
	}

	// Save caches
	_ = uc.enrichmentRepo.SaveSpellCache(spellCache)
	_ = uc.enrichmentRepo.SaveEquipmentCache(equipCache)

	logs = append(logs, fmt.Sprintf("Fetched: %d spells, %d equipment | Cached: %d spells, %d equipment",
		fetched["spells"], fetched["equipment"], cached["spells"], cached["equipment"]))

	// Create enriched character data
	charMap := uc.characterToMap(char)
	charMap["enriched"] = map[string]interface{}{
		"spells":    enrichedSpells,
		"equipment": enrichedEquipment,
	}

	// Write output
	outPath, err := uc.writeEnrichedCharacter(input.CharacterName, charMap, input.Inplace)
	if err != nil {
		return nil, err
	}

	logs = append(logs, fmt.Sprintf("Wrote enriched character to: %s (took %dms)", outPath, time.Since(start).Milliseconds()))

	return &EnrichCharacterOutput{
		OutputPath: outPath,
		Fetched:    fetched,
		Cached:     cached,
		Logs:       logs,
	}, nil
}

// FetchAllReferences fetches all spells and equipment from the API
func (uc *EnrichCharacterUseCase) FetchAllReferences(ctx context.Context) error {
	return uc.enrichmentRepo.FetchAllReferences(ctx)
}

// Helper methods

func (uc *EnrichCharacterUseCase) collectSpellNames(char *entities.Character) []string {
	set := make(map[string]struct{})
	for _, spell := range char.Spells {
		spell = strings.TrimSpace(spell)
		if spell != "" {
			set[spell] = struct{}{}
		}
	}
	for _, spell := range char.PreparedSpells {
		spell = strings.TrimSpace(spell)
		if spell != "" {
			set[spell] = struct{}{}
		}
	}
	
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (uc *EnrichCharacterUseCase) collectEquipmentNames(char *entities.Character) []string {
	set := make(map[string]struct{})
	
	if char.Weapon != "" {
		set[char.Weapon] = struct{}{}
	}
	if char.OffHand != "" {
		set[char.OffHand] = struct{}{}
	}
	if char.Armor != "" {
		set[char.Armor] = struct{}{}
	}
	if char.Shield != "" {
		set[char.Shield] = struct{}{}
	}
	for _, item := range char.Inventory {
		item = strings.TrimSpace(item)
		if item != "" {
			set[item] = struct{}{}
		}
	}
	
	names := make([]string, 0, len(set))
	for name := range set {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (uc *EnrichCharacterUseCase) sanitizeIndex(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, ",", "")
	return name
}

func (uc *EnrichCharacterUseCase) characterToMap(char *entities.Character) map[string]interface{} {
	data, _ := json.Marshal(char)
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	return result
}

func (uc *EnrichCharacterUseCase) writeEnrichedCharacter(name string, data map[string]interface{}, inplace bool) (string, error) {
	pretty, _ := json.MarshalIndent(data, "", "  ")
	
	if inplace {
		outPath := filepath.Join("data", "characters", name+".json")
		if err := os.WriteFile(outPath, pretty, 0644); err != nil {
			return "", fmt.Errorf("write %s: %w", outPath, err)
		}
		return outPath, nil
	}
	
	outDir := filepath.Join("data", "enrichments")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", outDir, err)
	}
	outPath := filepath.Join(outDir, name+".json")
	if err := os.WriteFile(outPath, pretty, 0644); err != nil {
		return "", fmt.Errorf("write %s: %w", outPath, err)
	}
	return outPath, nil
}

// fetchConcurrently fetches multiple items concurrently with rate limiting
func (uc *EnrichCharacterUseCase) fetchConcurrently(ctx context.Context, jobs []fetchJob, logs []string) []fetchResult {
	if len(jobs) == 0 {
		return nil
	}

	// Get rate limit from environment (requests per second)
	rateLimit := uc.getRateLimit()
	
	// Create channels
	resultsChan := make(chan fetchResult, len(jobs))
	semaphore := make(chan struct{}, 6) // Max 6 concurrent requests
	
	// Create ticker for rate limiting
	ticker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer ticker.Stop()

	var wg sync.WaitGroup

	// Launch workers
	for _, job := range jobs {
		wg.Add(1)
		go func(j fetchJob) {
			defer wg.Done()
			
			// Wait for rate limiter
			<-ticker.C
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Fetch the item
			result := fetchResult{
				kind:  j.kind,
				name:  j.name,
				index: j.index,
			}

			if j.kind == "spell" {
				info, err := uc.enrichmentRepo.FetchSpellInfo(ctx, j.index)
				result.spell = info
				result.err = err
			} else if j.kind == "equipment" {
				info, err := uc.enrichmentRepo.FetchEquipmentInfo(ctx, j.index)
				result.equip = info
				result.err = err
			}

			resultsChan <- result
		}(job)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var results []fetchResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// getRateLimit returns the rate limit from environment or default
func (uc *EnrichCharacterUseCase) getRateLimit() int {
	if s := os.Getenv("DND5EAPI_RPS"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			return v
		}
	}
	return 8 // Default: 8 requests per second
}

// fetchJob represents a fetch operation
type fetchJob struct {
	kind  string
	name  string
	index string
}

// fetchResult represents the result of a fetch operation
type fetchResult struct {
	kind  string
	name  string
	index string
	spell *repositories.SpellInfo
	equip *repositories.EquipmentInfo
	err   error
}

