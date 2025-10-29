package entities

import "testing"

func TestCharacterGetAbilityScore(t *testing.T) {
	char := &Character{
		Str: 16,
		Dex: 14,
		Con: 13,
		Int: 12,
		Wis: 10,
		Cha: 8,
	}

	tests := []struct {
		ability  string
		expected int
	}{
		{"str", 16},
		{"dex", 14},
		{"con", 13},
		{"int", 12},
		{"wis", 10},
		{"cha", 8},
		{"unknown", 10}, // default value
	}

	for _, tt := range tests {
		t.Run(tt.ability, func(t *testing.T) {
			result := char.GetAbilityScore(tt.ability)
			if result != tt.expected {
				t.Errorf("GetAbilityScore(%q) = %d; want %d", tt.ability, result, tt.expected)
			}
		})
	}
}

func TestCharacterHasSkillProficiency(t *testing.T) {
	char := &Character{
		Name:               "TestChar",
		SkillProficiencies: []string{"perception", "stealth", "athletics"},
	}

	tests := []struct {
		skill    string
		expected bool
	}{
		{"perception", true},
		{"stealth", true},
		{"athletics", true},
		{"acrobatics", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.skill, func(t *testing.T) {
			result := char.HasSkillProficiency(tt.skill)
			if result != tt.expected {
				t.Errorf("HasSkillProficiency(%q) = %v; want %v", tt.skill, result, tt.expected)
			}
		})
	}
}

func TestCharacterHasSkillProficiencyEmpty(t *testing.T) {
	char := &Character{
		Name:               "TestChar",
		SkillProficiencies: []string{},
	}

	if char.HasSkillProficiency("perception") {
		t.Error("Expected false for empty skill proficiencies")
	}
}
