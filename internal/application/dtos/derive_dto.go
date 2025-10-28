package dtos

// DerivedStatsDTO contains calculated character statistics
type DerivedStatsDTO struct {
	CharacterName string `json:"character_name"`
	
	// Ability Modifiers
	StrMod int `json:"str_mod"`
	DexMod int `json:"dex_mod"`
	ConMod int `json:"con_mod"`
	IntMod int `json:"int_mod"`
	WisMod int `json:"wis_mod"`
	ChaMod int `json:"cha_mod"`
	
	// Armor Class
	AC           int    `json:"ac"`
	ACCalculation string `json:"ac_calculation"`
	
	// Saving Throws
	StrSave        int    `json:"str_save"`
	DexSave        int    `json:"dex_save"`
	ConSave        int    `json:"con_save"`
	IntSave        int    `json:"int_save"`
	WisSave        int    `json:"wis_save"`
	ChaSave        int    `json:"cha_save"`
	SavingThrowStr string `json:"saving_throw_str"`
	
	// Skills
	Skills map[string]int `json:"skills"` // skill name -> modifier
	
	// Other Stats
	ProficiencyBonus int `json:"proficiency_bonus"`
	MaxHP            int `json:"max_hp,omitempty"`
}

// DeriveStatsRequestDTO represents a request to calculate derived stats
type DeriveStatsRequestDTO struct {
	CharacterName string `json:"character_name"`
}
