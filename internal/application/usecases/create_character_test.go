package usecases

import (
	"testing"

	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/entities"
)

const (
	errExpectedNilResult = "Expected nil result for invalid input"
)

// MockCharacterRepository is a simple mock for testing
type MockCharacterRepository struct {
	SaveFunc       func(character *entities.Character) error
	FindByNameFunc func(name string) (*entities.Character, error)
	DeleteFunc     func(name string) error
	ListFunc       func() ([]string, error)
	ExistsFunc     func(name string) bool
}

func (m *MockCharacterRepository) Save(character *entities.Character) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(character)
	}
	return nil
}

func (m *MockCharacterRepository) FindByName(name string) (*entities.Character, error) {
	if m.FindByNameFunc != nil {
		return m.FindByNameFunc(name)
	}
	return nil, nil
}

func (m *MockCharacterRepository) Delete(name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(name)
	}
	return nil
}

func (m *MockCharacterRepository) List() ([]string, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return []string{}, nil
}

func (m *MockCharacterRepository) Exists(name string) bool {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(name)
	}
	return false
}

func TestCreateCharacterSuccess(t *testing.T) {
	var savedCharacter *entities.Character
	mockRepo := &MockCharacterRepository{
		SaveFunc: func(character *entities.Character) error {
			savedCharacter = character
			return nil
		},
	}

	useCase := NewCreateCharacterUseCase(mockRepo)

	input := dtos.CreateCharacterDTO{
		Name:  "Gandalf",
		Race:  "Human",
		Class: "Wizard",
		Level: 5,
		Str:   10,
		Dex:   14,
		Con:   12,
		Int:   20,
		Wis:   16,
		Cha:   14,
	}

	result, err := useCase.Execute(input)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "Gandalf" {
		t.Errorf("Expected name 'Gandalf', got %s", result.Name)
	}

	if savedCharacter == nil {
		t.Fatal("Expected character to be saved")
	}

	if savedCharacter.Name != "Gandalf" {
		t.Errorf("Expected saved character name 'Gandalf', got %s", savedCharacter.Name)
	}
}

func TestCreateCharacterMissingName(t *testing.T) {
	mockRepo := &MockCharacterRepository{}
	useCase := NewCreateCharacterUseCase(mockRepo)

	input := dtos.CreateCharacterDTO{
		Race:  "Human",
		Class: "Wizard",
		Level: 5,
	}

	result, err := useCase.Execute(input)

	if err == nil {
		t.Error("Expected error for missing name, got nil")
	}

	if result != nil {
		t.Error(errExpectedNilResult)
	}
}

func TestCreateCharacterMissingRace(t *testing.T) {
	mockRepo := &MockCharacterRepository{}
	useCase := NewCreateCharacterUseCase(mockRepo)

	input := dtos.CreateCharacterDTO{
		Name:  "Gandalf",
		Class: "Wizard",
		Level: 5,
	}

	result, err := useCase.Execute(input)

	if err == nil {
		t.Error("Expected error for missing race, got nil")
	}

	if result != nil {
		t.Error(errExpectedNilResult)
	}
}

func TestCreateCharacterMissingClass(t *testing.T) {
	mockRepo := &MockCharacterRepository{}
	useCase := NewCreateCharacterUseCase(mockRepo)

	input := dtos.CreateCharacterDTO{
		Name:  "Gandalf",
		Race:  "Human",
		Level: 5,
	}

	result, err := useCase.Execute(input)

	if err == nil {
		t.Error("Expected error for missing class, got nil")
	}

	if result != nil {
		t.Error(errExpectedNilResult)
	}
}

func TestCreateCharacterRacialBonuses(t *testing.T) {
	mockRepo := &MockCharacterRepository{}
	useCase := NewCreateCharacterUseCase(mockRepo)

	tests := []struct {
		name        string
		race        string
		baseStr     int
		expectedStr int
	}{
		{"Human gets +1 to all", "Human", 10, 11},
		{"Dwarf gets +2 CON", "Dwarf", 10, 10},
		{"Elf no STR bonus", "Elf", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var saved *entities.Character
			mockRepo.SaveFunc = func(character *entities.Character) error {
				saved = character
				return nil
			}

			input := dtos.CreateCharacterDTO{
				Name:  "Test",
				Race:  tt.race,
				Class: "Fighter",
				Level: 1,
				Str:   tt.baseStr,
				Dex:   10,
				Con:   10,
				Int:   10,
				Wis:   10,
				Cha:   10,
			}

			_, err := useCase.Execute(input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if saved.Str != tt.expectedStr {
				t.Errorf("Expected STR %d, got %d", tt.expectedStr, saved.Str)
			}
		})
	}
}
