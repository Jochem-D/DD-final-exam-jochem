package usecases

import (
	"errors"
	"testing"

	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/services"
)

const (
	errMsgNotFound = "character not found"
	testWarrior    = "Warrior"
	testBarbarian  = "Barbarian"
)

// ========== DeriveCharacterStatsUseCase Tests ==========

func TestDeriveStatsSuccess(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:  testWarrior,
				Race:  "Human",
				Class: "Fighter",
				Level: 5,
				Str:   18,
				Dex:   14,
				Con:   16,
				Int:   10,
				Wis:   12,
				Cha:   8,
				Armor: "Chain Mail",
				SkillProficiencies: []string{"athletics", "intimidation"},
			}, nil
		},
	}

	useCase := NewDeriveCharacterStatsUseCase(mockRepo, services.NewCharacterService())
	stats, err := useCase.Execute(testWarrior)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if stats == nil {
		t.Fatal("Expected stats, got nil")
	}

	// Fighter level 5 has +3 proficiency bonus
	if stats.ProficiencyBonus != 3 {
		t.Errorf("Expected proficiency bonus 3, got %d", stats.ProficiencyBonus)
	}

	// STR 18 = +4 modifier
	if stats.StrMod != 4 {
		t.Errorf("Expected STR modifier 4, got %d", stats.StrMod)
	}

	// DEX 14 = +2 modifier
	if stats.DexMod != 2 {
		t.Errorf("Expected DEX modifier 2, got %d", stats.DexMod)
	}

	// Chain Mail AC = 16
	if stats.AC < 16 {
		t.Errorf("Expected AC >= 16 for Chain Mail, got %d", stats.AC)
	}
}

func TestDeriveStatsCharacterNotFound(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return nil, errors.New(errMsgNotFound)
		},
	}

	useCase := NewDeriveCharacterStatsUseCase(mockRepo, services.NewCharacterService())
	stats, err := useCase.Execute("NonExistent")

	if err == nil {
		t.Error("Expected error for non-existent character")
	}

	if stats != nil {
		t.Error("Expected nil stats for error case")
	}
}

// ========== ViewCharacterUseCase Tests ==========

func TestViewCharacterSuccess(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:               "Gandalf",
				Race:               "Human",
				Class:              "Wizard",
				Level:              10,
				Str:                10,
				Dex:                14,
				Con:                12,
				Int:                18,
				Wis:                16,
				Cha:                13,
				Background:         "Sage",
				SkillProficiencies: []string{"arcana", "history"},
				Spells:             []string{"Fireball", "Magic Missile"},
				Weapon:             "Staff",
			}, nil
		},
	}

	useCase := NewViewCharacterUseCase(mockRepo, services.NewCharacterService())
	output, err := useCase.Execute("Gandalf")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, got nil")
	}

	if output.Character.Name != "Gandalf" {
		t.Errorf("Expected name 'Gandalf', got %s", output.Character.Name)
	}

	if output.ArmorClass == 0 {
		t.Error("Expected AC to be calculated")
	}

	if output.InitiativeBonus == 0 {
		t.Error("Expected Initiative to be calculated")
	}

	if output.ProficiencyBonus == 0 {
		t.Error("Expected ProficiencyBonus to be calculated")
	}
}

func TestViewCharacterNotFound(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return nil, errors.New(errMsgNotFound)
		},
	}

	useCase := NewViewCharacterUseCase(mockRepo, services.NewCharacterService())
	output, err := useCase.Execute("NonExistent")

	if err == nil {
		t.Error("Expected error for non-existent character")
	}

	if output != nil {
		t.Error("Expected nil output for error case")
	}
}

// ========== GetLearnableSpellsUseCase Tests ==========

func TestGetLearnableSpellsSuccess(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:   "Sorcerer",
				Race:   "Human",
				Class:  "Sorcerer",
				Level:  3,
				Spells: []string{"Magic Missile"},
			}, nil
		},
	}

	mockSRD := &mockSRDRepository{
		GetLearnableSpellsFunc: func(className string, knownSpells []string) ([]string, error) {
			return []string{"Fireball", "Lightning Bolt", "Haste"}, nil
		},
	}

	useCase := NewGetLearnableSpellsUseCase(mockRepo, mockSRD)
	spells, err := useCase.Execute("Sorcerer")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if spells == nil || len(spells) == 0 {
		t.Error("Expected learnable spells, got none")
	}

	expectedCount := 3
	if len(spells) != expectedCount {
		t.Errorf("Expected %d learnable spells, got %d", expectedCount, len(spells))
	}
}

func TestGetLearnableSpellsNonSpellcaster(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:  testBarbarian,
				Race:  "Human",
				Class: "Barbarian",
				Level: 5,
			}, nil
		},
	}

	mockSRD := &mockSRDRepository{}

	useCase := NewGetLearnableSpellsUseCase(mockRepo, mockSRD)
	spells, err := useCase.Execute(testBarbarian)

	if err == nil {
		t.Error("Expected error for non-spellcaster class")
	}

	if spells != nil {
		t.Error("Expected nil spells for non-spellcaster")
	}
}

func TestGetLearnableSpellsCharacterNotFound(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return nil, errors.New(errMsgNotFound)
		},
	}

	mockSRD := &mockSRDRepository{}

	useCase := NewGetLearnableSpellsUseCase(mockRepo, mockSRD)
	spells, err := useCase.Execute("NonExistent")

	if err == nil {
		t.Error("Expected error for non-existent character")
	}

	if spells != nil {
		t.Error("Expected nil spells for error case")
	}
}

// ========== Additional Error Cases ==========

func TestLearnSpellInvalidClass(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:   testBarbarian,
				Class:  "Barbarian",
				Level:  5,
				Spells: []string{},
			}, nil
		},
	}

	mockSRD := &mockSRDRepository{
		IsValidSpellFunc:    func(s string) bool { return true },
		IsSpellForClassFunc: func(s, c string) (bool, error) { return true, nil },
	}

	useCase := NewLearnSpellUseCase(mockRepo, mockSRD)
	err := useCase.Execute(testBarbarian, "Fireball")

	if err == nil {
		t.Error("Expected error for non-spellcaster trying to learn spell")
	}
}

func TestPrepareSpellNotForClass(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:   "Cleric",
				Class:  "Cleric",
				Level:  5,
				Spells: []string{"Cure Wounds", "Guiding Bolt"},
			}, nil
		},
	}

	mockSRD := &mockSRDRepository{
		IsSpellForClassFunc: func(s, c string) (bool, error) {
			return false, nil // Spell not for this class
		},
		GetSpellLevelFunc: func(s string) (int, error) {
			return 1, nil
		},
	}

	useCase := NewPrepareSpellUseCase(mockRepo, mockSRD)

	input := PrepareSpellInput{
		CharacterName: "Cleric",
		SpellName:     "Eldritch Blast", // Warlock spell
	}

	err := useCase.Execute(input)

	if err == nil {
		t.Error("Expected error when preparing spell not available to class")
	}
}
