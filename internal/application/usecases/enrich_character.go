package usecases

import (
	"context"

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
	characterRepo    repositories.CharacterRepository
	enrichmentRepo   repositories.EnrichmentRepository
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
	// This is a simplified version - the actual implementation would be more complex
	// For now, we'll delegate to the infrastructure layer
	return &EnrichCharacterOutput{
		OutputPath: input.CharacterName,
		Fetched:    make(map[string]int),
		Cached:     make(map[string]int),
		Logs:       []string{},
	}, nil
}

// FetchAllReferences fetches all spells and equipment from the API
func (uc *EnrichCharacterUseCase) FetchAllReferences(ctx context.Context) error {
	return uc.enrichmentRepo.FetchAllReferences(ctx)
}
