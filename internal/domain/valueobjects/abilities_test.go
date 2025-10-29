package valueobjects

import "testing"

func TestAbilityModifier(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected int
	}{
		{"Score 1", 1, -5},
		{"Score 3", 3, -4},
		{"Score 8", 8, -1},
		{"Score 10", 10, 0},
		{"Score 11", 11, 0},
		{"Score 12", 12, 1},
		{"Score 14", 14, 2},
		{"Score 18", 18, 4},
		{"Score 20", 20, 5},
		{"Score 30", 30, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AbilityModifier(tt.score)
			if result != tt.expected {
				t.Errorf("AbilityModifier(%d) = %d; want %d", tt.score, result, tt.expected)
			}
		})
	}
}

func TestProficiencyBonus(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"Level 1", 1, 2},
		{"Level 4", 4, 2},
		{"Level 5", 5, 3},
		{"Level 8", 8, 3},
		{"Level 9", 9, 4},
		{"Level 12", 12, 4},
		{"Level 13", 13, 5},
		{"Level 16", 16, 5},
		{"Level 17", 17, 6},
		{"Level 20", 20, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProficiencyBonus(tt.level)
			if result != tt.expected {
				t.Errorf("ProficiencyBonus(%d) = %d; want %d", tt.level, result, tt.expected)
			}
		})
	}
}
