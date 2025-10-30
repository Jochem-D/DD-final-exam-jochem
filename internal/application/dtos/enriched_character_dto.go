package dtos

// EnrichedCharacterDTO contains all character data plus all calculated/derived stats
// This is meant to be a complete representation for the frontend
type EnrichedCharacterDTO struct {
	// Base character data
	CharacterDTO
	
	// Ability Modifiers (calculated)
	StrMod int `json:"str_mod"`
	DexMod int `json:"dex_mod"`
	ConMod int `json:"con_mod"`
	IntMod int `json:"int_mod"`
	WisMod int `json:"wis_mod"`
	ChaMod int `json:"cha_mod"`
	
	// Derived Stats
	ProficiencyBonus  int    `json:"proficiency_bonus"`
	ArmorClass        int    `json:"armor_class"`
	ACCalculation     string `json:"ac_calculation"`
	Initiative        int    `json:"initiative"`
	PassivePerception int    `json:"passive_perception"`
	Speed             int    `json:"speed"`
	MaxHP             int    `json:"max_hp"` 
	
	// Saving Throws (with proficiency applied)
	StrSave int  `json:"str_save"`
	DexSave int  `json:"dex_save"`
	ConSave int  `json:"con_save"`
	IntSave int  `json:"int_save"`
	WisSave int  `json:"wis_save"`
	ChaSave int  `json:"cha_save"`
	
	// Saving Throw Proficiencies (which saves the class is proficient in)
	StrSaveProf bool `json:"str_save_prof"`
	DexSaveProf bool `json:"dex_save_prof"`
	ConSaveProf bool `json:"con_save_prof"`
	IntSaveProf bool `json:"int_save_prof"`
	WisSaveProf bool `json:"wis_save_prof"`
	ChaSaveProf bool `json:"cha_save_prof"`
	
	// Skills (calculated modifiers with proficiency)
	Skills map[string]int `json:"skills"` // skill name -> total modifier
	
	// Skill Proficiencies (which skills character is proficient in)
	SkillProfs map[string]bool `json:"skill_profs"` // skill name -> is proficient
	
	// Weapon Attacks (calculated with ability mods and proficiency)
	WeaponAttacks []WeaponAttackDTO `json:"weapon_attacks,omitempty"`
}

// WeaponAttackDTO represents a calculated weapon attack
type WeaponAttackDTO struct {
	Name        string `json:"name"`
	AttackBonus string `json:"attack_bonus"` // e.g., "+5"
	Damage      string `json:"damage"`       // e.g., "1d8+3 slashing"
}

