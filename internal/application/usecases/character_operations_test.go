package usecases

import (
	"errors"
	"testing"

	"ddsheetfinal/internal/domain/entities"
)

const errExpectedNoError = "Expected no error, got %v"

func TestDeleteCharacterSuccess(t *testing.T) {
	deletedName := ""
	mockRepo := &MockCharacterRepository{
		DeleteFunc: func(name string) error {
			deletedName = name
			return nil
		},
	}

	useCase := NewDeleteCharacterUseCase(mockRepo)
	err := useCase.Execute("Gandalf")

	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}

	if deletedName != "Gandalf" {
		t.Errorf("Expected to delete 'Gandalf', got %s", deletedName)
	}
}

func TestDeleteCharacterError(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		DeleteFunc: func(name string) error {
			return errors.New("character not found")
		},
	}

	useCase := NewDeleteCharacterUseCase(mockRepo)
	err := useCase.Execute("NonExistent")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestListCharactersSuccess(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		ListFunc: func() ([]string, error) {
			return []string{"Gandalf", "Aragorn", "Frodo"}, nil
		},
	}

	useCase := NewListCharactersUseCase(mockRepo)
	result, err := useCase.Execute()

	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}

	if len(result.Names) != 3 {
		t.Errorf("Expected 3 characters, got %d", len(result.Names))
	}

	if result.Names[0] != "Gandalf" {
		t.Errorf("Expected first character 'Gandalf', got %s", result.Names[0])
	}
}

func TestListCharactersEmpty(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		ListFunc: func() ([]string, error) {
			return []string{}, nil
		},
	}

	useCase := NewListCharactersUseCase(mockRepo)
	result, err := useCase.Execute()

	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}

	if len(result.Names) != 0 {
		t.Errorf("Expected 0 characters, got %d", len(result.Names))
	}
}

func TestGetCharacterSuccess(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return &entities.Character{
				Name:  name,
				Race:  "Human",
				Class: "Wizard",
				Level: 5,
			}, nil
		},
	}

	useCase := NewGetCharacterUseCase(mockRepo)
	result, err := useCase.Execute("Gandalf")

	if err != nil {
		t.Errorf(errExpectedNoError, err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Name != "Gandalf" {
		t.Errorf("Expected name 'Gandalf', got %s", result.Name)
	}
}

func TestGetCharacterNotFound(t *testing.T) {
	mockRepo := &MockCharacterRepository{
		FindByNameFunc: func(name string) (*entities.Character, error) {
			return nil, errors.New("character not found")
		},
	}

	useCase := NewGetCharacterUseCase(mockRepo)
	result, err := useCase.Execute("NonExistent")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result")
	}
}
