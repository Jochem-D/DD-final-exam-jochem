package valueobjects

import "testing"

func TestAbilityModifierEdgeCases(t *testing.T) {
	tests := []struct {
		score    int
		expected int
	}{
		{1, -5},   // Minimum possible
		{3, -4},
		{8, -1},
		{9, -1},
		{10, 0},   // Average
		{11, 0},
		{12, 1},
		{18, 4},
		{20, 5},   // Common max
		{30, 10},  // Legendary
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := AbilityModifier(tt.score)
			if result != tt.expected {
				t.Errorf("AbilityModifier(%d) = %d; want %d", tt.score, result, tt.expected)
			}
		})
	}
}

func TestProficiencyBonusAllLevels(t *testing.T) {
	tests := []struct {
		level    int
		expected int
	}{
		{1, 2},
		{2, 2},
		{3, 2},
		{4, 2},
		{5, 3},
		{8, 3},
		{9, 4},
		{12, 4},
		{13, 5},
		{16, 5},
		{17, 6},
		{20, 6},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := ProficiencyBonus(tt.level)
			if result != tt.expected {
				t.Errorf("ProficiencyBonus(%d) = %d; want %d", tt.level, result, tt.expected)
			}
		})
	}
}

func TestIsSpellcasterAllClasses(t *testing.T) {
	spellcasters := []string{"wizard", "cleric", "druid", "bard", "sorcerer", "warlock", "paladin", "ranger"}
	
	for _, class := range spellcasters {
		if !IsSpellcaster(class) {
			t.Errorf("IsSpellcaster(%q) should be true", class)
		}
	}
}

func TestIsSpellcasterNonCasters(t *testing.T) {
	nonCasters := []string{"fighter", "barbarian", "rogue", "monk"}
	
	for _, class := range nonCasters {
		if IsSpellcaster(class) {
			t.Errorf("IsSpellcaster(%q) should be false", class)
		}
	}
}

func TestIsPreparedCasterAllClasses(t *testing.T) {
	tests := []struct {
		class    string
		expected bool
	}{
		{"wizard", true},
		{"cleric", true},
		{"druid", true},
		{"paladin", true},
		{"bard", false},
		{"sorcerer", false},
		{"warlock", false},
		{"ranger", false},
		{"fighter", false},
	}

	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			result := IsPreparedCaster(tt.class)
			if result != tt.expected {
				t.Errorf("IsPreparedCaster(%q) = %v; want %v", tt.class, result, tt.expected)
			}
		})
	}
}

func TestGetArmorInfoAllTypes(t *testing.T) {
	armors := []struct {
		name     string
		shouldExist bool
	}{
		{"Leather Armor", true},
		{"Chain Mail", true},
		{"Plate Armor", true},
		{"Studded Leather", true},
		{"Hide Armor", true},
		{"Nonexistent Armor", false},
		{"", false},
	}

	for _, tt := range armors {
		t.Run(tt.name, func(t *testing.T) {
			info, exists := GetArmorInfo(tt.name)
			if exists != tt.shouldExist {
				t.Errorf("GetArmorInfo(%q) exists = %v; want %v", tt.name, exists, tt.shouldExist)
			}
			if exists && info.Base == 0 {
				t.Errorf("GetArmorInfo(%q) should have non-zero base AC", tt.name)
			}
		})
	}
}

func TestAllClassesHaveSavingThrows(t *testing.T) {
	classes := []string{"barbarian", "bard", "cleric", "druid", "fighter", "monk", "paladin", "ranger", "rogue", "sorcerer", "warlock", "wizard"}
	
	for _, class := range classes {
		saves, exists := ClassSavingThrowProficiencies[class]
		if !exists {
			t.Errorf("Class %q missing from ClassSavingThrowProficiencies", class)
			continue
		}
		if len(saves) != 2 {
			t.Errorf("Class %q should have exactly 2 saving throw proficiencies, got %d", class, len(saves))
		}
	}
}

func TestAllSkillsHaveAbilityMapping(t *testing.T) {
	requiredSkills := []string{
		"acrobatics", "animal handling", "arcana", "athletics",
		"deception", "history", "insight", "intimidation",
		"investigation", "medicine", "nature", "perception",
		"performance", "persuasion", "religion", "sleight of hand",
		"stealth", "survival",
	}

	for _, skill := range requiredSkills {
		ability, exists := SkillToAbility[skill]
		if !exists {
			t.Errorf("Skill %q missing from SkillToAbility map", skill)
			continue
		}
		validAbilities := []string{"str", "dex", "con", "int", "wis", "cha"}
		found := false
		for _, valid := range validAbilities {
			if ability == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Skill %q has invalid ability %q", skill, ability)
		}
	}
}

func TestGetSpellSlotsProgression(t *testing.T) {
	// Test wizard spell slot progression across multiple levels
	levels := []struct {
		level        int
		expectedL1   int
		expectedL2   int
		expectedL3   int
	}{
		{1, 2, 0, 0},
		{3, 4, 2, 0},
		{5, 4, 3, 2},
		{9, 4, 3, 3},
	}

	for _, tt := range levels {
		slots := GetSpellSlots("wizard", tt.level)
		if slots == nil {
			t.Fatalf("GetSpellSlots(wizard, %d) returned nil", tt.level)
		}

		if slots[1] != tt.expectedL1 {
			t.Errorf("Wizard level %d: expected %d 1st-level slots, got %d", tt.level, tt.expectedL1, slots[1])
		}
		if slots[2] != tt.expectedL2 {
			t.Errorf("Wizard level %d: expected %d 2nd-level slots, got %d", tt.level, tt.expectedL2, slots[2])
		}
		if slots[3] != tt.expectedL3 {
			t.Errorf("Wizard level %d: expected %d 3rd-level slots, got %d", tt.level, tt.expectedL3, slots[3])
		}
	}
}

func TestBackgroundSkillsComplete(t *testing.T) {
	backgrounds := []string{
		"acolyte", "charlatan", "criminal", "entertainer",
		"folk hero", "guild artisan", "hermit", "noble",
		"outlander", "sage", "sailor", "soldier", "urchin",
	}

	for _, bg := range backgrounds {
		skills, exists := BackgroundSkillProficiencies[bg]
		if !exists {
			t.Errorf("Background %q missing from BackgroundSkillProficiencies", bg)
			continue
		}
		if len(skills) != 2 {
			t.Errorf("Background %q should have exactly 2 skill proficiencies, got %d", bg, len(skills))
		}
	}
}

func TestClassSkillChoicesAllClasses(t *testing.T) {
	classes := []string{"barbarian", "bard", "cleric", "druid", "fighter", "monk", "paladin", "ranger", "rogue", "sorcerer", "warlock", "wizard"}
	
	for _, class := range classes {
		choices, exists := ClassSkillChoices[class]
		if !exists {
			t.Errorf("Class %q missing from ClassSkillChoices", class)
			continue
		}
		if choices < 2 || choices > 4 {
			t.Errorf("Class %q has invalid skill choices: %d (should be 2-4)", class, choices)
		}
	}
}

func TestRacialBonusesCoverage(t *testing.T) {
	races := []string{"dwarf", "elf", "halfling", "human", "dragonborn", "gnome", "half-elf", "half-orc", "tiefling"}
	
	for _, race := range races {
		bonuses, exists := RacialAbilityBonuses[race]
		if !exists {
			t.Errorf("Race %q missing from RacialAbilityBonuses", race)
			continue
		}
		if len(bonuses) == 0 {
			t.Errorf("Race %q has no ability bonuses", race)
		}
	}
}
