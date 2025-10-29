package usecases

import (
	"errors"
	"testing"

	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/entities"
)

const (
	errNoError        = "Expected no error, got %v"
	errCharSaved      = "Expected character to be saved"
	spellCureWounds   = "Cure Wounds"
)

func TestEquipItemSuccess(t *testing.T) {
	var saved *entities.Character
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:      "Fighter",
				Race:      "Human",
				Class:     "Fighter",
				Level:     5,
				Inventory: []string{"Longsword", "Chain Mail"},
			}, nil
		},
		SaveFunc: func(character *entities.Character) error {
			saved = character
			return nil
		},
	}

	// simple SRD mock
	mockSRD := &mockSRDRepository{
		IsValidEquipmentFunc: func(e string) bool { return true },
	}

	useCase := NewEquipItemUseCase(mockRepo, mockSRD)

	input := EquipItemInput{
		CharacterName: "Fighter",
		Weapon:        "Longsword",
		Slot:          "main hand",
	}

	err := useCase.Execute(input)

	if err != nil {
		t.Errorf(errNoError, err)
	}

	if saved == nil {
		t.Fatal(errCharSaved)
	}

	if saved.Weapon != "Longsword" {
		t.Errorf("Expected weapon 'Longsword', got %s", saved.Weapon)
	}
}

func TestEquipItemCharacterNotFound(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return nil, errors.New("character not found")
		},
	}

	mockSRD := &mockSRDRepository{IsValidEquipmentFunc: func(e string) bool { return true }}
	useCase := NewEquipItemUseCase(mockRepo, mockSRD)

	input := EquipItemInput{
		CharacterName: "NonExistent",
		Weapon:        "Sword",
		Slot:          "main hand",
	}

	err := useCase.Execute(input)

	if err == nil {
		t.Error("Expected error for non-existent character")
	}
}

func TestUnequipItemSuccess(t *testing.T) {
	var saved *entities.Character
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:   "Fighter",
				Race:   "Human",
				Class:  "Fighter",
				Level:  5,
				Weapon: "Longsword",
			}, nil
		},
		SaveFunc: func(character *entities.Character) error {
			saved = character
			return nil
		},
	}

	useCase := NewUnequipItemUseCase(mockRepo)

	input := UnequipItemInput{
		CharacterName: "Fighter",
		UnequipWeapon: true,
	}

	err := useCase.Execute(input)

	if err != nil {
		t.Errorf(errNoError, err)
	}

	if saved == nil {
		t.Fatal(errCharSaved)
	}

	if saved.Weapon != "" {
		t.Errorf("Expected weapon to be empty, got %s", saved.Weapon)
	}
}

func TestLearnSpellSuccess(t *testing.T) {
	var saved *entities.Character
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:   "Bard",
				Race:   "Human",
				Class:  "Bard",
				Level:  5,
				Spells: []string{"Magic Missile"},
			}, nil
		},
		SaveFunc: func(character *entities.Character) error {
			saved = character
			return nil
		},
	}

	mockSRD := &mockSRDRepository{
		IsValidSpellFunc:    func(s string) bool { return true },
		IsSpellForClassFunc: func(s, c string) (bool, error) { return true, nil },
		GetSpellLevelFunc:   func(s string) (int, error) { return 3, nil },
	}
	useCase := NewLearnSpellUseCase(mockRepo, mockSRD)

	err := useCase.Execute("Bard", "Fireball")

	if err != nil {
		t.Errorf(errNoError, err)
	}

	if saved == nil {
		t.Fatal(errCharSaved)
	}

	found := false
	for _, spell := range saved.Spells {
		if spell == "Fireball" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected Fireball to be in learned spells")
	}
}

func TestPrepareSpellSuccess(t *testing.T) {
	var saved *entities.Character
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:   "Cleric",
				Race:   "Human",
				Class:  "Cleric",
				Level:  5,
				Spells: []string{spellCureWounds, "Bless"},
			}, nil
		},
		SaveFunc: func(character *entities.Character) error {
			saved = character
			return nil
		},
	}

	mockSRD := &mockSRDRepository{
		IsSpellForClassFunc: func(s, c string) (bool, error) { return true, nil },
		GetSpellLevelFunc:   func(s string) (int, error) { return 1, nil },
	}
	useCase := NewPrepareSpellUseCase(mockRepo, mockSRD)

	input := PrepareSpellInput{
		CharacterName: "Cleric",
		SpellName:     spellCureWounds,
	}

	err := useCase.Execute(input)

	if err != nil {
		t.Errorf(errNoError, err)
	}

	if saved == nil {
		t.Fatal(errCharSaved)
	}

	found := false
	for _, spell := range saved.PreparedSpells {
		if spell == spellCureWounds {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected Cure Wounds to be in prepared spells")
	}
}

func TestSaveCharacterSuccess(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		SaveFunc: func(character *entities.Character) error {
			return nil
		},
	}

	useCase := NewSaveCharacterUseCase(mockRepo)

	char := &dtos.CharacterDTO{
		Name:  "TestChar",
		Race:  "Human",
		Class: "Fighter",
		Level: 1,
	}

	err := useCase.Execute(char)

	if err != nil {
		t.Errorf(errNoError, err)
	}
}

func TestSaveCharacterError(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		SaveFunc: func(character *entities.Character) error {
			return errors.New("save failed")
		},
	}

	useCase := NewSaveCharacterUseCase(mockRepo)

	char := &dtos.CharacterDTO{
		Name:  "TestChar",
		Race:  "Human",
		Class: "Fighter",
		Level: 1,
	}

	err := useCase.Execute(char)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
