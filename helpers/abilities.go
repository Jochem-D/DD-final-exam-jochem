// helpers/abilities.go
package helpers

import "math"

// AbilityMod converts an ability score to its 5e modifier.
// Use floor((score-10)/2) so negative modifiers round down correctly (e.g. 9 -> -1).
func AbilityMod(score int) int {
	return int(math.Floor(float64(score-10) / 2.0))
}

// GetProficiencyBonus returns the 5e proficiency bonus for a character level.
func GetProficiencyBonus(level int) int {
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
