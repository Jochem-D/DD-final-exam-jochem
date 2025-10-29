package services

import (
	"testing"

	"ddsheetfinal/internal/domain/entities"
)

func TestComputeArmorClass(t *testing.T) {
	service := NewCharacterService()

	tests := []struct {
		name        string
		char        *entities.Character
		expectedAC  int
		description string
	}{
		{
			name: "Unarmored character with 14 DEX",
			char: &entities.Character{
				Name:  "Test",
				Class: "Wizard",
				Dex:   14,
			},
			expectedAC:  12, // 10 + 2 (DEX mod)
			description: "10 + Dex(+2)",
		},
		{
			name: "Monk unarmored (DEX 16, WIS 18)",
			char: &entities.Character{
				Name:  "Monk",
				Class: "Monk",
				Dex:   16,
				Wis:   18,
			},
			expectedAC:  17, // 10 + 3 (DEX) + 4 (WIS)
			description: "10 + Dex(+3) + Wis(+4)",
		},
		{
			name: "Barbarian unarmored (DEX 14, CON 16)",
			char: &entities.Character{
				Name:  "Barbarian",
				Class: "Barbarian",
				Dex:   14,
				Con:   16,
			},
			expectedAC:  15, // 10 + 2 (DEX) + 3 (CON)
			description: "10 + Dex(+2) + Con(+3)",
		},
		{
			name: "Light armor (leather) with DEX 18",
			char: &entities.Character{
				Name:  "Rogue",
				Class: "Rogue",
				Dex:   18,
				Armor: "Leather",
			},
			expectedAC:  15, // 11 (base) + 4 (DEX)
			description: "11 (light base) + Dex(+4)",
		},
		{
			name: "Medium armor (breastplate) with DEX 16",
			char: &entities.Character{
				Name:  "Fighter",
				Class: "Fighter",
				Dex:   16,
				Armor: "Breastplate",
			},
			expectedAC:  16, // 14 (base) + 2 (DEX capped)
			description: "14 (medium base) + Dex cap(+2)",
		},
		{
			name: "Heavy armor (plate) ignores DEX",
			char: &entities.Character{
				Name:  "Paladin",
				Class: "Paladin",
				Dex:   10,
				Armor: "Plate",
			},
			expectedAC:  18, // 18 (base) + 0 (no DEX)
			description: "18 (heavy base)",
		},
		{
			name: "Shield adds +2",
			char: &entities.Character{
				Name:   "Fighter",
				Class:  "Fighter",
				Dex:    14,
				Armor:  "Chain Mail",
				Shield: "Shield",
			},
			expectedAC:  18, // 16 (heavy base) + 2 (shield)
			description: "16 (heavy base) + shield(+2)",
		},
		{
			name: "Unarmored with shield",
			char: &entities.Character{
				Name:   "Cleric",
				Class:  "Cleric",
				Dex:    12,
				Shield: "Shield",
			},
			expectedAC:  13, // 10 + 1 (DEX) + 2 (shield)
			description: "10 + Dex(+1) + shield(+2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, desc := service.ComputeArmorClass(tt.char)
			if ac != tt.expectedAC {
				t.Errorf("ComputeArmorClass() AC = %d; want %d", ac, tt.expectedAC)
			}
			if desc != tt.description {
				t.Errorf("ComputeArmorClass() description = %q; want %q", desc, tt.description)
			}
		})
	}
}

func TestComputeInitiativeBonus(t *testing.T) {
	service := NewCharacterService()

	tests := []struct {
		name     string
		dex      int
		expected int
	}{
		{"DEX 10", 10, 0},
		{"DEX 14", 14, 2},
		{"DEX 18", 18, 4},
		{"DEX 8", 8, -1},
		{"DEX 20", 20, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &entities.Character{Dex: tt.dex}
			result := service.ComputeInitiativeBonus(char)
			if result != tt.expected {
				t.Errorf("ComputeInitiativeBonus() = %d; want %d", result, tt.expected)
			}
		})
	}
}

func TestComputePassivePerception(t *testing.T) {
	service := NewCharacterService()

	tests := []struct {
		name       string
		wis        int
		level      int
		proficient bool
		expected   int
	}{
		{"WIS 10, not proficient", 10, 1, false, 10},
		{"WIS 14, not proficient", 14, 1, false, 12},
		{"WIS 14, proficient level 1", 14, 1, true, 14},  // 10 + 2 (WIS) + 2 (prof)
		{"WIS 18, proficient level 5", 18, 5, true, 17},  // 10 + 4 (WIS) + 3 (prof)
		{"WIS 20, proficient level 17", 20, 17, true, 21}, // 10 + 5 (WIS) + 6 (prof)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &entities.Character{
				Wis:   tt.wis,
				Level: tt.level,
			}
			if tt.proficient {
				char.SkillProficiencies = []string{"perception"}
			}
			
			result := service.ComputePassivePerception(char)
			if result != tt.expected {
				t.Errorf("ComputePassivePerception() = %d; want %d", result, tt.expected)
			}
		})
	}
}

func TestIsSaveProficient(t *testing.T) {
	service := NewCharacterService()

	tests := []struct {
		name     string
		class    string
		ability  string
		expected bool
	}{
		// Wizard: Int, Wis
		{"Wizard INT save", "Wizard", "int", true},
		{"Wizard WIS save", "Wizard", "wis", true},
		{"Wizard STR save", "Wizard", "str", false},
		
		// Fighter: Str, Con
		{"Fighter STR save", "Fighter", "str", true},
		{"Fighter CON save", "Fighter", "con", true},
		{"Fighter DEX save", "Fighter", "dex", false},
		
		// Case insensitive
		{"Upper case class", "WIZARD", "int", true},
		{"Full ability name", "Wizard", "intelligence", true},
		
		// Unknown class
		{"Unknown class", "UnknownClass", "str", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &entities.Character{Class: tt.class}
			result := service.IsSaveProficient(char, tt.ability)
			if result != tt.expected {
				t.Errorf("IsSaveProficient(%q, %q) = %v; want %v", tt.class, tt.ability, result, tt.expected)
			}
		})
	}
}

func TestComputeSkillModifier(t *testing.T) {
	service := NewCharacterService()

	tests := []struct {
		name       string
		skill      string
		ability    int
		level      int
		proficient bool
		expected   int
	}{
		{"Perception (WIS 14), not proficient", "perception", 14, 1, false, 2},
		{"Perception (WIS 14), proficient", "perception", 14, 1, true, 4}, // 2 (WIS) + 2 (prof)
		{"Stealth (DEX 18), proficient level 5", "stealth", 18, 5, true, 7}, // 4 (DEX) + 3 (prof)
		{"Athletics (STR 16), proficient level 9", "athletics", 16, 9, true, 7}, // 3 (STR) + 4 (prof)
		{"Investigation (INT 20), proficient level 17", "investigation", 20, 17, true, 11}, // 5 (INT) + 6 (prof)
		{"Unknown skill", "unknown", 14, 1, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := &entities.Character{
				Level: tt.level,
			}
			
			// Set ability score based on skill
			skillLower := tt.skill
			switch skillLower {
			case "perception":
				char.Wis = tt.ability
			case "stealth":
				char.Dex = tt.ability
			case "athletics":
				char.Str = tt.ability
			case "investigation":
				char.Int = tt.ability
			}
			
			if tt.proficient {
				char.SkillProficiencies = []string{tt.skill}
			}
			
			result := service.ComputeSkillModifier(char, tt.skill)
			if result != tt.expected {
				t.Errorf("ComputeSkillModifier(%q) = %d; want %d", tt.skill, result, tt.expected)
			}
		})
	}
}
