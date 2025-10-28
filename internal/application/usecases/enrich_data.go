package usecases

import (
	"context"
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/repositories"
	"strings"
)

// EnrichDataUseCase handles fetching enriched data from external API
type EnrichDataUseCase struct {
	enrichmentRepo repositories.EnrichmentRepository
}

// NewEnrichDataUseCase creates a new EnrichDataUseCase
func NewEnrichDataUseCase(enrichmentRepo repositories.EnrichmentRepository) *EnrichDataUseCase {
	return &EnrichDataUseCase{
		enrichmentRepo: enrichmentRepo,
	}
}

// ExecuteSpell fetches spell information from external API
func (uc *EnrichDataUseCase) ExecuteSpell(ctx context.Context, spellIndex string) (*dtos.EnrichedSpellDTO, error) {
	spellInfo, err := uc.enrichmentRepo.FetchSpellInfo(ctx, spellIndex)
	if err != nil {
		return nil, err
	}

	return &dtos.EnrichedSpellDTO{
		Name:        spellInfo.Name,
		School:      spellInfo.School,
		Range:       spellInfo.Range,
		Description: spellInfo.Description,
	}, nil
}

// ExecuteEquipment fetches equipment information from external API
func (uc *EnrichDataUseCase) ExecuteEquipment(ctx context.Context, equipIndex string) (*dtos.EnrichedEquipmentDTO, error) {
	equipInfo, err := uc.enrichmentRepo.FetchEquipmentInfo(ctx, equipIndex)
	if err != nil {
		return nil, err
	}

	return &dtos.EnrichedEquipmentDTO{
		Name:                equipInfo.Name,
		Category:            equipInfo.Category,
		RangeNormal:         equipInfo.RangeNormal,
		TwoHanded:           equipInfo.TwoHanded,
		ArmorBase:           equipInfo.ArmorBase,
		DexAllowed:          equipInfo.DexAllowed,
		MaxDex:              equipInfo.MaxDex,
		StrMinimum:          equipInfo.StrMinimum,
		StealthDisadvantage: equipInfo.StealthDisadvantage,
		PriceGP:             equipInfo.PriceGP,
		DamageText:          equipInfo.DamageText,
		DamageAvg:           equipInfo.DamageAvg,
	}, nil
}

// SanitizeIndex converts a name to an API index format
func (uc *EnrichDataUseCase) SanitizeIndex(name string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
}
