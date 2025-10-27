package usecases

import (
	"ddsheetfinal/internal/domain/repositories"
)

// DeleteCharacterUseCase handles character deletion
type DeleteCharacterUseCase struct {
	characterRepo repositories.CharacterRepository
}

// NewDeleteCharacterUseCase creates a new DeleteCharacterUseCase
func NewDeleteCharacterUseCase(characterRepo repositories.CharacterRepository) *DeleteCharacterUseCase {
	return &DeleteCharacterUseCase{
		characterRepo: characterRepo,
	}
}

// Execute deletes a character by name
func (uc *DeleteCharacterUseCase) Execute(name string) error {
	return uc.characterRepo.Delete(name)
}
