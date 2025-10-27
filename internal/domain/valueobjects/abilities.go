package valueobjects

import "math"

// AbilityModifier calculates the D&D 5e ability modifier for a given score
func AbilityModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2.0))
}

// ProficiencyBonus returns the proficiency bonus for a given character level
func ProficiencyBonus(level int) int {
	switch {
	case level >= 17:
		return 6
	case level >= 13:
		return 5
	case level >= 9:
		return 4
	case level >= 5:
		return 3
	default:
		return 2
	}
}
