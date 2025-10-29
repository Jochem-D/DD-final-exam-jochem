package valueobjects

import "testing"

func TestSkillToAbilityMapping(t *testing.T) {
	tests := []struct {
		skill    string
		expected string
	}{
		{"athletics", "str"},
		{"acrobatics", "dex"},
		{"stealth", "dex"},
		{"arcana", "int"},
		{"history", "int"},
		{"perception", "wis"},
		{"insight", "wis"},
		{"persuasion", "cha"},
		{"deception", "cha"},
	}

	for _, tt := range tests {
		t.Run(tt.skill, func(t *testing.T) {
			ability, exists := SkillToAbility[tt.skill]
			if !exists {
				t.Errorf("Skill %q not found in SkillToAbility map", tt.skill)
				return
			}
			if ability != tt.expected {
				t.Errorf("SkillToAbility[%q] = %q; want %q", tt.skill, ability, tt.expected)
			}
		})
	}
}

func TestClassSavingThrowProficiencies(t *testing.T) {
	tests := []struct {
		class     string
		abilities []string
	}{
		{"fighter", []string{"str", "con"}},
		{"wizard", []string{"int", "wis"}},
		{"cleric", []string{"wis", "cha"}},
		{"rogue", []string{"dex", "int"}},
	}

	for _, tt := range tests {
		t.Run(tt.class, func(t *testing.T) {
			saves, exists := ClassSavingThrowProficiencies[tt.class]
			if !exists {
				t.Errorf("Class %q not found in ClassSavingThrowProficiencies", tt.class)
				return
			}
			if len(saves) != len(tt.abilities) {
				t.Errorf("Expected %d saving throws, got %d", len(tt.abilities), len(saves))
				return
			}
			for i, ability := range tt.abilities {
				if saves[i] != ability {
					t.Errorf("Expected saving throw %d to be %q, got %q", i, ability, saves[i])
				}
			}
		})
	}
}

func TestRacialBonuses(t *testing.T) {
	tests := []struct {
		race string
		attr string
		has  bool
	}{
		{"human", "str", true},
		{"human", "dex", true},
		{"dwarf", "con", true},
		{"elf", "dex", true},
		{"halfling", "dex", true},
	}

	for _, tt := range tests {
		t.Run(tt.race+"_"+tt.attr, func(t *testing.T) {
			bonuses, exists := RacialAbilityBonuses[tt.race]
			if !exists && tt.has {
				t.Errorf("Race %q not found in RacialAbilityBonuses", tt.race)
				return
			}
			if exists {
				_, hasBonus := bonuses[tt.attr]
				if hasBonus != tt.has {
					t.Errorf("Race %q should have %q bonus: %v, got: %v", tt.race, tt.attr, tt.has, hasBonus)
				}
			}
		})
	}
}

func TestClassAndBackgroundSkillMappings(t *testing.T) {
	// Class skills - rogue should offer stealth and sleight of hand
	rogueSkills, ok := ClassSkillProficiencies["rogue"]
	if !ok {
		t.Fatal("rogue not present in ClassSkillProficiencies")
	}
	foundStealth := false
	foundSleight := false
	for _, s := range rogueSkills {
		if s == "stealth" {
			foundStealth = true
		}
		if s == "sleight of hand" {
			foundSleight = true
		}
	}
	if !foundStealth {
		t.Error("rogue skills should include stealth")
	}

	if !foundSleight {
		t.Error("rogue skills should include sleight of hand")
	}

	// Background skills - sage should include arcana and history
	sageSkills, ok := BackgroundSkillProficiencies["sage"]
	if !ok {
		t.Fatal("sage not present in BackgroundSkillProficiencies")
	}
	want := map[string]bool{"arcana": false, "history": false}
	for _, s := range sageSkills {
		if _, exists := want[s]; exists {
			want[s] = true
		}
	}
	for k, v := range want {
		if !v {
			t.Errorf("background 'sage' should include %s", k)
		}
	}
}

func TestGetSpellSlotsFullCaster(t *testing.T) {
	// Wizard level 1: should have cantrips and 1st-level slots
	slots := GetSpellSlots("wizard", 1)
	if slots == nil {
		t.Fatal("expected spell slots for wizard level 1")
	}
	if slots[0] != 3 { // wizard level1 returns 3 cantrips in this simplified model
		t.Errorf("wizard level1 cantrips = %d; want 3", slots[0])
	}
	if slots[1] != 2 {
		t.Errorf("wizard level1 1st-level slots = %d; want 2", slots[1])
	}

	// Bard level 5: check cantrips and spell slots copied from FullCasterSpellSlots
	slots = GetSpellSlots("bard", 5)
	if slots == nil {
		t.Fatal("expected spell slots for bard level 5")
	}
	if slots[0] != 3 { // bard at level 5 gets 3 cantrips per getCantripsKnown
		t.Errorf("bard level5 cantrips = %d; want 3", slots[0])
	}
	if slots[3] != 0 && slots[3] < 2 { // ensure higher levels present per table
		t.Errorf("unexpected bard level5 4th-level slots: %v", slots)
	}
}

func TestGetSpellSlotsHalfCaster(t *testing.T) {
	// Paladin is a half-caster - no cantrips, but should have 1st-level slots at level 2
	slots := GetSpellSlots("paladin", 2)
	if slots == nil {
		t.Fatal("expected spell slots for paladin level 2")
	}
	if slots[0] != 0 {
		t.Errorf("paladin should not have cantrips; got %d", slots[0])
	}
	if slots[1] != 2 {
		t.Errorf("paladin level2 1st-level slots = %d; want 2", slots[1])
	}
}

func TestGetSpellSlotsWarlock(t *testing.T) {
	// Warlock has pact magic: check cantrips and single-level filled
	slots := GetSpellSlots("warlock", 1)
	if slots == nil {
		t.Fatal("expected spell slots for warlock level 1")
	}
	if slots[0] != 2 {
		t.Errorf("warlock level1 cantrips = %d; want 2", slots[0])
	}
	// At level 1 warlock has 1 slot at some level; ensure total slots sum > 0
	sum := 0
	for _, v := range slots {
		sum += v
	}
	if sum <= 2 { // only cantrips would be 2, so sum must be >2
		t.Errorf("warlock level1 total slots sum = %d; want >2", sum)
	}
}

func TestGetSpellSlotsInvalidLevel(t *testing.T) {
	if GetSpellSlots("wizard", 0) != nil {
		t.Error("expected nil for level 0")
	}
	if GetSpellSlots("wizard", 21) != nil {
		t.Error("expected nil for level 21")
	}
}
