package repositories

// EquipmentInfo contains enriched equipment data
type EquipmentInfo struct {
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

// SpellInfo contains enriched spell data
type SpellInfo struct {
	Name        string `json:"name"`
	School      string `json:"school"`
	Range       string `json:"range"`
	Description string `json:"description,omitempty"`
}

// SRDRepository defines the contract for accessing SRD data
type SRDRepository interface {
	// IsValidSpell checks if a spell exists in the SRD
	IsValidSpell(spellName string) bool
	
	// IsValidEquipment checks if equipment exists in the SRD
	IsValidEquipment(equipmentName string) bool
	
	// IsSpellForClass checks if a spell is available to a class
	IsSpellForClass(spellName, className string) (bool, error)
	
	// GetLearnableSpells returns spells a character can learn
	GetLearnableSpells(className string, knownSpells []string) ([]string, error)
	
	// GetSpellLevel returns the level of a spell (0-9)
	GetSpellLevel(spellName string) (int, error)
}
