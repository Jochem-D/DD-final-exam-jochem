package usecases

import (
	"fmt"
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/repositories"
)

// SaveCharacterUseCase handles saving/updating a character
type SaveCharacterUseCase struct {
	characterRepo repositories.CharacterRepository
}

// NewSaveCharacterUseCase creates a new SaveCharacterUseCase
func NewSaveCharacterUseCase(characterRepo repositories.CharacterRepository) *SaveCharacterUseCase {
	return &SaveCharacterUseCase{
		characterRepo: characterRepo,
	}
}

// Execute saves or updates a character
func (uc *SaveCharacterUseCase) Execute(characterDTO *dtos.CharacterDTO) error {
	if characterDTO == nil {
		return fmt.Errorf("character data is required")
	}
	
	if characterDTO.Name == "" {
		return fmt.Errorf("character name is required")
	}

	// Convert DTO to entity
	character := dtos.ToCharacterEntity(characterDTO)

	// Save character
	return uc.characterRepo.Save(character)
}
