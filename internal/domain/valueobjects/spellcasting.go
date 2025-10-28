package valueobjects

// Spellcasting class constants to avoid duplication across use cases

// SpellcastingClasses returns all classes that can cast spells
var SpellcastingClasses = []string{"bard", "cleric", "druid", "paladin", "ranger", "sorcerer", "warlock", "wizard"}

// PreparedCasterClasses returns classes that prepare spells (vs learning them)
var PreparedCasterClasses = []string{"cleric", "druid", "paladin", "wizard"}

// IsSpellcaster checks if a class can cast spells
func IsSpellcaster(class string) bool {
	for _, c := range SpellcastingClasses {
		if c == class {
			return true
		}
	}
	return false
}

// IsPreparedCaster checks if a class prepares spells
func IsPreparedCaster(class string) bool {
	for _, c := range PreparedCasterClasses {
		if c == class {
			return true
		}
	}
	return false
}

// GetSpellcastingAbility returns the primary spellcasting ability for a class
func GetSpellcastingAbility(class string) string {
	switch class {
	case "wizard":
		return "int"
	case "cleric", "druid", "ranger":
		return "wis"
	case "bard", "paladin", "sorcerer", "warlock":
		return "cha"
	default:
		return ""
	}
}
