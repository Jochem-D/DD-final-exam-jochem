package dtos

// EnrichRequestDTO represents a request to enrich character data from external API
type EnrichRequestDTO struct {
	Type  string `json:"type"`  // "spell" or "equipment"
	Index string `json:"index"` // API index (e.g., "fireball", "longsword")
}

// EnrichedSpellDTO represents enriched spell information
type EnrichedSpellDTO struct {
	Name        string `json:"name"`
	School      string `json:"school"`
	Range       string `json:"range"`
	Description string `json:"description,omitempty"`
}

// EnrichedEquipmentDTO represents enriched equipment information
type EnrichedEquipmentDTO struct {
	Name                string  `json:"name"`
	Category            string  `json:"category"`
	RangeNormal         int     `json:"range_normal,omitempty"`
	TwoHanded           bool    `json:"two_handed,omitempty"`
	ArmorBase           int     `json:"armor_base,omitempty"`
	DexAllowed          bool    `json:"dex_allowed,omitempty"`
	MaxDex              *int    `json:"max_dex,omitempty"`
	StrMinimum          int     `json:"str_minimum,omitempty"`
	StealthDisadvantage bool    `json:"stealth_disadvantage,omitempty"`
	PriceGP             float64 `json:"price_gp,omitempty"`
	DamageText          string  `json:"damage_text,omitempty"`
	DamageAvg           float64 `json:"damage_avg,omitempty"`
}

// EnrichCharacterRequestDTO represents a request to enrich all character data
type EnrichCharacterRequestDTO struct {
	CharacterName string `json:"character_name"`
}

// EnrichCharacterResponseDTO represents the result of enriching character data
type EnrichCharacterResponseDTO struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
