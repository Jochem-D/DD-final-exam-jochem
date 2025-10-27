package repositories

import "context"

// EnrichmentRepository defines the contract for fetching and caching enrichment data
type EnrichmentRepository interface {
	// FetchSpellInfo retrieves spell information from the external API
	FetchSpellInfo(ctx context.Context, spellIndex string) (*SpellInfo, error)
	
	// FetchEquipmentInfo retrieves equipment information from the external API
	FetchEquipmentInfo(ctx context.Context, equipmentIndex string) (*EquipmentInfo, error)
	
	// LoadSpellCache loads the spell cache
	LoadSpellCache() (map[string]SpellInfo, error)
	
	// SaveSpellCache saves the spell cache
	SaveSpellCache(cache map[string]SpellInfo) error
	
	// LoadEquipmentCache loads the equipment cache
	LoadEquipmentCache() (map[string]EquipmentInfo, error)
	
	// SaveEquipmentCache saves the equipment cache
	SaveEquipmentCache(cache map[string]EquipmentInfo) error
	
	// FetchAllReferences fetches all spells and equipment from the API
	FetchAllReferences(ctx context.Context) error
}
