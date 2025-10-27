package external

import (
	"context"
	"ddsheetfinal/helpers"
	"ddsheetfinal/internal/domain/repositories"
)

// DND5EAPIEnrichmentRepository implements EnrichmentRepository using the dnd5eapi
type DND5EAPIEnrichmentRepository struct {
	cacheDir string
}

// NewDND5EAPIEnrichmentRepository creates a new enrichment repository
func NewDND5EAPIEnrichmentRepository(cacheDir string) repositories.EnrichmentRepository {
	return &DND5EAPIEnrichmentRepository{
		cacheDir: cacheDir,
	}
}

// FetchSpellInfo retrieves spell information from the external API
func (r *DND5EAPIEnrichmentRepository) FetchSpellInfo(ctx context.Context, spellIndex string) (*repositories.SpellInfo, error) {
	// Delegate to the existing helpers for now
	// In a full refactor, we'd move the API logic here
	logs := &helpers.StdoutLog{}
	info, err := helpers.FetchSpellInfoPublic(ctx, spellIndex, logs)
	if err != nil {
		return nil, err
	}
	
	return &repositories.SpellInfo{
		Name:        info.Name,
		School:      info.School,
		Range:       info.Range,
		Description: info.Description,
	}, nil
}

// FetchEquipmentInfo retrieves equipment information from the external API
func (r *DND5EAPIEnrichmentRepository) FetchEquipmentInfo(ctx context.Context, equipmentIndex string) (*repositories.EquipmentInfo, error) {
	// Delegate to the existing helpers for now
	logs := &helpers.StdoutLog{}
	info, err := helpers.FetchEquipmentInfoPublic(ctx, equipmentIndex, logs)
	if err != nil {
		return nil, err
	}
	
	return &repositories.EquipmentInfo{
		Name:                info.Name,
		Category:            info.Category,
		RangeNormal:         info.RangeNormal,
		TwoHanded:           info.TwoHanded,
		ArmorBase:           info.ArmorBase,
		DexAllowed:          info.DexAllowed,
		MaxDex:              info.MaxDex,
		StrMinimum:          info.StrMinimum,
		StealthDisadvantage: info.StealthDisadvantage,
		PriceGP:             info.PriceGP,
		DamageText:          info.DamageText,
		DamageAvg:           info.DamageAvg,
	}, nil
}

// LoadSpellCache loads the spell cache
func (r *DND5EAPIEnrichmentRepository) LoadSpellCache() (map[string]repositories.SpellInfo, error) {
	cache, err := helpers.LoadSpellCache()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]repositories.SpellInfo)
	for k, v := range cache {
		result[k] = repositories.SpellInfo{
			Name:        v.Name,
			School:      v.School,
			Range:       v.Range,
			Description: v.Description,
		}
	}
	return result, nil
}

// SaveSpellCache saves the spell cache
func (r *DND5EAPIEnrichmentRepository) SaveSpellCache(cache map[string]repositories.SpellInfo) error {
	helperCache := make(map[string]helpers.SpellInfo)
	for k, v := range cache {
		helperCache[k] = helpers.SpellInfo{
			Name:        v.Name,
			School:      v.School,
			Range:       v.Range,
			Description: v.Description,
		}
	}
	return helpers.SaveSpellCache(helperCache)
}

// LoadEquipmentCache loads the equipment cache
func (r *DND5EAPIEnrichmentRepository) LoadEquipmentCache() (map[string]repositories.EquipmentInfo, error) {
	cache, err := helpers.LoadEquipCache()
	if err != nil {
		return nil, err
	}
	
	result := make(map[string]repositories.EquipmentInfo)
	for k, v := range cache {
		result[k] = repositories.EquipmentInfo{
			Name:                v.Name,
			Category:            v.Category,
			RangeNormal:         v.RangeNormal,
			TwoHanded:           v.TwoHanded,
			ArmorBase:           v.ArmorBase,
			DexAllowed:          v.DexAllowed,
			MaxDex:              v.MaxDex,
			StrMinimum:          v.StrMinimum,
			StealthDisadvantage: v.StealthDisadvantage,
			PriceGP:             v.PriceGP,
			DamageText:          v.DamageText,
			DamageAvg:           v.DamageAvg,
		}
	}
	return result, nil
}

// SaveEquipmentCache saves the equipment cache
func (r *DND5EAPIEnrichmentRepository) SaveEquipmentCache(cache map[string]repositories.EquipmentInfo) error {
	helperCache := make(map[string]helpers.EquipmentInfo)
	for k, v := range cache {
		helperCache[k] = helpers.EquipmentInfo{
			Name:                v.Name,
			Category:            v.Category,
			RangeNormal:         v.RangeNormal,
			TwoHanded:           v.TwoHanded,
			ArmorBase:           v.ArmorBase,
			DexAllowed:          v.DexAllowed,
			MaxDex:              v.MaxDex,
			StrMinimum:          v.StrMinimum,
			StealthDisadvantage: v.StealthDisadvantage,
			PriceGP:             v.PriceGP,
			DamageText:          v.DamageText,
			DamageAvg:           v.DamageAvg,
		}
	}
	return helpers.SaveEquipCache(helperCache)
}

// FetchAllReferences fetches all spells and equipment from the API
func (r *DND5EAPIEnrichmentRepository) FetchAllReferences(ctx context.Context) error {
	logs := &helpers.StdoutLog{}
	return helpers.FetchAllReferences(ctx, logs)
}
