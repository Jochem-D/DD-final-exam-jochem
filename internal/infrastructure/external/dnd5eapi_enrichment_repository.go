package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ddsheetfinal/internal/domain/repositories"
)

const (
	apiBaseLegacy = "https://www.dnd5eapi.co/api"
	apiBase2014   = "https://www.dnd5eapi.co/api/2014"
)

// DND5EAPIEnrichmentRepository implements EnrichmentRepository using the dnd5eapi
type DND5EAPIEnrichmentRepository struct {
	cacheDir   string
	httpClient *http.Client
}

// NewDND5EAPIEnrichmentRepository creates a new enrichment repository
func NewDND5EAPIEnrichmentRepository(cacheDir string) repositories.EnrichmentRepository {
	// Ensure cache directory exists
	_ = os.MkdirAll(cacheDir, 0o755)
	
	return &DND5EAPIEnrichmentRepository{
		cacheDir: cacheDir,
		httpClient: &http.Client{
			Timeout: 12 * time.Second,
		},
	}
}

// FetchSpellInfo retrieves spell information from the external API or cache
func (r *DND5EAPIEnrichmentRepository) FetchSpellInfo(ctx context.Context, spellIndex string) (*repositories.SpellInfo, error) {
	// Try cache first
	cache, _ := r.LoadSpellCache()
	if cache == nil {
		cache = make(map[string]repositories.SpellInfo)
	}
	if info, ok := cache[spellIndex]; ok {
		return &info, nil
	}
	
	// Fetch from API
	var body struct {
		Name   string `json:"name"`
		Range  string `json:"range"`
		School struct {
			Name string `json:"name"`
		} `json:"school"`
		Desc []string `json:"desc"`
	}
	
	path := "/spells/" + spellIndex
	if err := r.getJSONWithFallback(ctx, path, &body); err != nil {
		return nil, fmt.Errorf("fetch spell %s: %w", spellIndex, err)
	}
	
	desc := strings.Join(body.Desc, "\n\n")
	info := &repositories.SpellInfo{
		Name:        body.Name,
		School:      body.School.Name,
		Range:       body.Range,
		Description: desc,
	}
	
	// Save to cache
	cache[spellIndex] = *info
	_ = r.SaveSpellCache(cache)
	
	return info, nil
}

// FetchEquipmentInfo retrieves equipment information from the external API
// FetchEquipmentInfo retrieves equipment information from the external API or cache
func (r *DND5EAPIEnrichmentRepository) FetchEquipmentInfo(ctx context.Context, equipmentIndex string) (*repositories.EquipmentInfo, error) {
	// Try cache first
	cache, _ := r.LoadEquipmentCache()
	if cache == nil {
		cache = make(map[string]repositories.EquipmentInfo)
	}
	if info, ok := cache[equipmentIndex]; ok {
		return &info, nil
	}
	
	// Fetch from API
	var raw map[string]any
	
	path := "/equipment/" + equipmentIndex
	if err := r.getJSONWithFallback(ctx, path, &raw); err != nil {
		return nil, fmt.Errorf("fetch equipment %s: %w", equipmentIndex, err)
	}
	
	info := &repositories.EquipmentInfo{}
	
	// Parse basic fields
	if v, ok := raw["name"].(string); ok {
		info.Name = v
	}
	if v, ok := raw["equipment_category"].(map[string]any); ok {
		if n, ok2 := v["name"].(string); ok2 {
			info.Category = n
		}
	}
	
	// Parse weapon fields
	if v, ok := raw["range"].(map[string]any); ok {
		if normal, ok2 := v["normal"].(float64); ok2 {
			info.RangeNormal = int(normal)
		}
	}
	if v, ok := raw["two_handed_damage"].(map[string]any); ok {
		info.TwoHanded = v != nil
	}
	
	// Parse armor fields
	if v, ok := raw["armor_class"].(map[string]any); ok {
		if base, ok2 := v["base"].(float64); ok2 {
			info.ArmorBase = int(base)
		}
		if dex, ok2 := v["dex_bonus"].(bool); ok2 {
			info.DexAllowed = dex
		}
		if maxDex, ok2 := v["max_bonus"].(float64); ok2 {
			maxDexInt := int(maxDex)
			info.MaxDex = &maxDexInt
		}
	}
	
	if sd, ok := raw["stealth_disadvantage"].(bool); ok {
		info.StealthDisadvantage = sd
	}
	if sm, ok := raw["str_minimum"].(float64); ok {
		info.StrMinimum = int(sm)
	}
	
	// Parse cost
	if cost, ok := raw["cost"].(map[string]any); ok {
		if qty, ok2 := cost["quantity"].(float64); ok2 {
			if unit, ok3 := cost["unit"].(string); ok3 {
				info.PriceGP = r.convertToGP(qty, unit)
			}
		}
	}
	
	// Parse damage
	if dmg, ok := raw["damage"].(map[string]any); ok {
		if dice, ok2 := dmg["damage_dice"].(string); ok2 {
			info.DamageText = dice
		}
	}
	
	// Save to cache
	cache[equipmentIndex] = *info
	_ = r.SaveEquipmentCache(cache)
	
	return info, nil
}

// LoadSpellCache loads the spell cache
func (r *DND5EAPIEnrichmentRepository) LoadSpellCache() (map[string]repositories.SpellInfo, error) {
	path := filepath.Join(r.cacheDir, "spells.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]repositories.SpellInfo), nil
		}
		return nil, err
	}
	
	var cache map[string]repositories.SpellInfo
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	
	return cache, nil
}

// SaveSpellCache saves the spell cache
func (r *DND5EAPIEnrichmentRepository) SaveSpellCache(cache map[string]repositories.SpellInfo) error {
	path := filepath.Join(r.cacheDir, "spells.json")
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadEquipmentCache loads the equipment cache
func (r *DND5EAPIEnrichmentRepository) LoadEquipmentCache() (map[string]repositories.EquipmentInfo, error) {
	path := filepath.Join(r.cacheDir, "equipment.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]repositories.EquipmentInfo), nil
		}
		return nil, err
	}
	
	var cache map[string]repositories.EquipmentInfo
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	
	return cache, nil
}

// SaveEquipmentCache saves the equipment cache
func (r *DND5EAPIEnrichmentRepository) SaveEquipmentCache(cache map[string]repositories.EquipmentInfo) error {
	path := filepath.Join(r.cacheDir, "equipment.json")
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// FetchAllReferences fetches all spells and equipment from the API
func (r *DND5EAPIEnrichmentRepository) FetchAllReferences(ctx context.Context) error {
	// Fetch all spells
	fmt.Println("Fetching spell list...")
	var spellList struct {
		Results []struct {
			Index string `json:"index"`
			Name  string `json:"name"`
		} `json:"results"`
	}
	
	if err := r.getJSONWithFallback(ctx, "/spells", &spellList); err != nil {
		return fmt.Errorf("fetch spell list: %w", err)
	}
	
	fmt.Printf("Fetching %d spells with concurrency...\n", len(spellList.Results))
	spellCache := r.fetchSpellsConcurrently(ctx, spellList.Results)
	
	if err := r.SaveSpellCache(spellCache); err != nil {
		return fmt.Errorf("save spell cache: %w", err)
	}
	fmt.Printf("✓ Cached %d spells\n", len(spellCache))
	
	// Fetch all equipment
	fmt.Println("Fetching equipment list...")
	var equipList struct {
		Results []struct {
			Index string `json:"index"`
			Name  string `json:"name"`
		} `json:"results"`
	}
	
	if err := r.getJSONWithFallback(ctx, "/equipment", &equipList); err != nil {
		return fmt.Errorf("fetch equipment list: %w", err)
	}
	
	fmt.Printf("Fetching %d equipment items with concurrency...\n", len(equipList.Results))
	equipCache := r.fetchEquipmentConcurrently(ctx, equipList.Results)
	
	if err := r.SaveEquipmentCache(equipCache); err != nil {
		return fmt.Errorf("save equipment cache: %w", err)
	}
	fmt.Printf("✓ Cached %d equipment items\n", len(equipCache))
	
	return nil
}

// fetchSpellsConcurrently fetches spells concurrently using goroutines
func (r *DND5EAPIEnrichmentRepository) fetchSpellsConcurrently(ctx context.Context, spells []struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}) map[string]repositories.SpellInfo {
	const maxWorkers = 10
	type result struct {
		index string
		info  *repositories.SpellInfo
		err   error
	}
	
	results := make(chan result, len(spells))
	semaphore := make(chan struct{}, maxWorkers)
	
	for _, spell := range spells {
		spell := spell // capture loop variable
		go func() {
			semaphore <- struct{}{} // acquire
			defer func() { <-semaphore }() // release
			
			info, err := r.FetchSpellInfo(ctx, spell.Index)
			results <- result{index: spell.Index, info: info, err: err}
		}()
	}
	
	cache := make(map[string]repositories.SpellInfo)
	for i := 0; i < len(spells); i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("  Warning: failed to fetch spell %s: %v\n", res.index, res.err)
			continue
		}
		if res.info != nil {
			cache[res.index] = *res.info
		}
		if (i+1)%50 == 0 {
			fmt.Printf("  Progress: %d/%d spells\n", i+1, len(spells))
		}
	}
	
	return cache
}

// fetchEquipmentConcurrently fetches equipment concurrently using goroutines
func (r *DND5EAPIEnrichmentRepository) fetchEquipmentConcurrently(ctx context.Context, equipment []struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}) map[string]repositories.EquipmentInfo {
	const maxWorkers = 10
	type result struct {
		index string
		info  *repositories.EquipmentInfo
		err   error
	}
	
	results := make(chan result, len(equipment))
	semaphore := make(chan struct{}, maxWorkers)
	
	for _, equip := range equipment {
		equip := equip // capture loop variable
		go func() {
			semaphore <- struct{}{} // acquire
			defer func() { <-semaphore }() // release
			
			info, err := r.FetchEquipmentInfo(ctx, equip.Index)
			results <- result{index: equip.Index, info: info, err: err}
		}()
	}
	
	cache := make(map[string]repositories.EquipmentInfo)
	for i := 0; i < len(equipment); i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("  Warning: failed to fetch equipment %s: %v\n", res.index, res.err)
			continue
		}
		if res.info != nil {
			cache[res.index] = *res.info
		}
		if (i+1)%50 == 0 {
			fmt.Printf("  Progress: %d/%d equipment\n", i+1, len(equipment))
		}
	}
	
	return cache
}

// Helper methods

func (r *DND5EAPIEnrichmentRepository) getJSONWithFallback(ctx context.Context, path string, dst any) error {
	// Try legacy API first
	u1 := apiBaseLegacy + path
	if err := r.getJSON(ctx, u1, dst); err == nil {
		return nil
	}
	
	// Try 2014 API
	u2 := apiBase2014 + path
	return r.getJSON(ctx, u2, dst)
}

func (r *DND5EAPIEnrichmentRepository) getJSON(ctx context.Context, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "DDsheetfinal/1.0")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	
	return json.NewDecoder(resp.Body).Decode(dst)
}

func (r *DND5EAPIEnrichmentRepository) convertToGP(qty float64, unit string) float64 {
	switch strings.ToLower(unit) {
	case "cp":
		return qty / 100
	case "sp":
		return qty / 10
	case "gp":
		return qty
	case "pp":
		return qty * 10
	default:
		return qty
	}
}

// LoadCharacterEnrichment loads enriched data for a specific character from the enrichments directory
func (r *DND5EAPIEnrichmentRepository) LoadCharacterEnrichment(characterName string) (map[string]interface{}, error) {
	enrichmentPath := filepath.Join("data", "enrichments", characterName+".json")
	
	data, err := os.ReadFile(enrichmentPath)
	if err != nil {
		return nil, fmt.Errorf("read enrichment file for %s: %w", characterName, err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse enrichment file for %s: %w", characterName, err)
	}
	
	return result, nil
}
