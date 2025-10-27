package entities

// Character represents a D&D 5e character with all their attributes and equipment
type Character struct {
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

// GetAbilityScore returns the ability score for a given ability abbreviation
func (c *Character) GetAbilityScore(ability string) int {
	switch ability {
	case "str":
		return c.Str
	case "dex":
		return c.Dex
	case "con":
		return c.Con
	case "int":
		return c.Int
	case "wis":
		return c.Wis
	case "cha":
		return c.Cha
	default:
		return 10
	}
}

// HasSkillProficiency checks if the character has proficiency in a skill
func (c *Character) HasSkillProficiency(skill string) bool {
	for _, s := range c.SkillProficiencies {
		if s == skill {
			return true
		}
	}
	return false
}
