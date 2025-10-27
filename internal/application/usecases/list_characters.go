package usecases

import (
	"ddsheetfinal/internal/domain/repositories"
)

// ListCharactersUseCase handles listing all characters
type ListCharactersUseCase struct {
	characterRepo repositories.CharacterRepository
}

// NewListCharactersUseCase creates a new ListCharactersUseCase
func NewListCharactersUseCase(characterRepo repositories.CharacterRepository) *ListCharactersUseCase {
	return &ListCharactersUseCase{
		characterRepo: characterRepo,
	}
}

// Execute returns a list of all character names
func (uc *ListCharactersUseCase) Execute() ([]string, error) {
	return uc.characterRepo.List()
}
