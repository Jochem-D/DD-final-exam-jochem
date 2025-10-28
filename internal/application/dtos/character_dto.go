package dtos

// CharacterDTO represents a character for the presentation layer
// This decouples the presentation from domain entities
type CharacterDTO struct {
	Name       string   `json:"name"`
	Race       string   `json:"race"`
	Class      string   `json:"class"`
	Background string   `json:"background,omitempty"`
	Level      int      `json:"level"`
	Str        int      `json:"str"`
	Dex        int      `json:"dex"`
	Con        int      `json:"con"`
	Int        int      `json:"int"`
	Wis        int      `json:"wis"`
	Cha        int      `json:"cha"`

	ProficiencyBonus   int      `json:"proficiency_bonus,omitempty"`
	SkillProficiencies []string `json:"skill_proficiencies,omitempty"`

	Weapon  string `json:"weapon,omitempty"`
	OffHand string `json:"off_hand,omitempty"`
	Armor   string `json:"armor,omitempty"`
	Shield  string `json:"shield,omitempty"`

	Inventory []string `json:"inventory,omitempty"`

	Spells         []string `json:"spells,omitempty"`
	PreparedSpells []string `json:"prepared_spells,omitempty"`
}

// CreateCharacterDTO represents input for creating a character
type CreateCharacterDTO struct {
	Name       string   `json:"name"`
	Race       string   `json:"race"`
	Class      string   `json:"class"`
	Background string   `json:"background,omitempty"`
	Level      int      `json:"level"`
	Str        int      `json:"str"`
	Dex        int      `json:"dex"`
	Con        int      `json:"con"`
	Int        int      `json:"int"`
	Wis        int      `json:"wis"`
	Cha        int      `json:"cha"`
	Skills     []string `json:"skills,omitempty"`
}

// CharacterListDTO represents a list of character names
type CharacterListDTO struct {
	Names []string `json:"names"`
}

// EquipItemDTO represents input for equipping an item
type EquipItemDTO struct {
	CharacterName string `json:"character_name"`
	ItemName      string `json:"item_name"`
	Slot          string `json:"slot"` // weapon, off-hand, armor, shield
}

// UnequipItemDTO represents input for unequipping an item
type UnequipItemDTO struct {
	CharacterName string `json:"character_name"`
	Slot          string `json:"slot"`
}

// SpellDTO represents spell information
type SpellDTO struct {
	Name        string `json:"name"`
	School      string `json:"school"`
	Range       string `json:"range"`
	Description string `json:"description,omitempty"`
}

// LearnSpellDTO represents input for learning a spell
type LearnSpellDTO struct {
	CharacterName string `json:"character_name"`
	SpellName     string `json:"spell_name"`
}

// PrepareSpellDTO represents input for preparing a spell
type PrepareSpellDTO struct {
	CharacterName string `json:"character_name"`
	SpellName     string `json:"spell_name"`
}

// LearnableSpellsDTO represents learnable spells for a character
type LearnableSpellsDTO struct {
	CharacterName string   `json:"character_name"`
	Spells        []string `json:"spells"`
}
