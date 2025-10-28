package usecases

import (
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/repositories"
)

// GetCharacterUseCase retrieves a character without derived stats
type GetCharacterUseCase struct {
	characterRepo repositories.CharacterRepository
}

// NewGetCharacterUseCase creates a new GetCharacterUseCase
func NewGetCharacterUseCase(characterRepo repositories.CharacterRepository) *GetCharacterUseCase {
	return &GetCharacterUseCase{
		characterRepo: characterRepo,
	}
}

// Execute retrieves a character by name
func (uc *GetCharacterUseCase) Execute(name string) (*dtos.CharacterDTO, error) {
	character, err := uc.characterRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	return dtos.ToCharacterDTO(character), nil
}
