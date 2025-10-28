package cli

import (
	"context"
	"ddsheetfinal/internal/application/usecases"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EnrichHandler handles the enrich command
type EnrichHandler struct {
	enrichUseCase *usecases.EnrichCharacterUseCase
}

// NewEnrichHandler creates a new enrich handler
func NewEnrichHandler(enrichUseCase *usecases.EnrichCharacterUseCase) *EnrichHandler {
	return &EnrichHandler{
		enrichUseCase: enrichUseCase,
	}
}

// Handle executes the enrich command
func (h *EnrichHandler) Handle(args []string) {
	inplace, force, fetchAll, name := h.parseArgs(args)

	baseCtx := context.Background()

	// If --fetch-all is specified, pre-cache all references
	if fetchAll {
		if err := h.runFetchAll(baseCtx); err != nil {
			fmt.Printf("fetch-all failed: %v\n", err)
			os.Exit(1)
		}
		// If no name was provided, we ran in fetch-only mode: exit successfully
		if name == "" {
			fmt.Println("✅ fetch-only completed")
			return
		}
	}

	// Ensure we have a character name
	if name == "" {
		fmt.Println("Error: character name required")
		fmt.Println("Usage: enrich <name> [--inplace] [--force] [--fetch-all]")
		os.Exit(1)
	}

	// 60s is reasonable for a single-character enrich run
	ctx, cancel := context.WithTimeout(baseCtx, 60*time.Second)
	defer cancel()

	input := usecases.EnrichCharacterInput{
		CharacterName: name,
		Inplace:       inplace,
		Force:         force,
	}

	out, err := h.enrichUseCase.Execute(ctx, input)
	if err != nil {
		fmt.Printf("❌ Enrich failed: %v\n", err)
		os.Exit(1)
	}

	// Display results
	fmt.Printf("\n✅ Enrichment Complete!\n")
	fmt.Printf("Output: %s\n", out.OutputPath)
	if len(out.Fetched) > 0 {
		fmt.Printf("Fetched: %v\n", out.Fetched)
	}
	if len(out.Cached) > 0 {
		fmt.Printf("Cached: %v\n", out.Cached)
	}
	if len(out.Logs) > 0 {
		fmt.Println("\nLogs:")
		for _, log := range out.Logs {
			fmt.Printf("  %s\n", log)
		}
	}
}

// parseArgs parses flags and the optional name argument from args
func (h *EnrichHandler) parseArgs(args []string) (inplace, force, fetchAll bool, name string) {
	fs := flag.NewFlagSet("enrich", flag.ExitOnError)
	
	// Allow flags before or after the name by registering them first
	inplacePtr := fs.Bool("inplace", false, "update the original JSON instead of writing to data/enrichments/<name>.json")
	forcePtr := fs.Bool("force", false, "force fetching/caching even if enriched file exists")
	fetchAllPtr := fs.Bool("fetch-all", false, "fetch all spells and equipment from the API into the cache before enriching")
	
	// Parse all args - flag package handles flags anywhere in args
	_ = fs.Parse(args)

	// Name is the first non-flag argument
	remaining := fs.Args()
	if len(remaining) > 0 {
		name = strings.TrimSpace(remaining[0])
	}

	return *inplacePtr, *forcePtr, *fetchAllPtr, name
}

// runFetchAll performs the long-running fetch-all operation under a timeout
func (h *EnrichHandler) runFetchAll(baseCtx context.Context) error {
	fmt.Println("🔄 fetch-all: fetching all spells and equipment into cache (this may take a while)...")
	ctxFetch, cancelFetch := context.WithTimeout(baseCtx, 10*time.Minute)
	defer cancelFetch()
	
	// Use the enrichment repository to fetch all references
	if err := h.enrichUseCase.FetchAllReferences(ctxFetch); err != nil {
		return err
	}
	
	fmt.Println("fetch-all: done")
	return nil
}

// Helper to save enriched character
func (h *EnrichHandler) saveEnrichedCharacter(name string, data map[string]interface{}, inplace bool) (string, error) {
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
