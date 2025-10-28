package usecases

import (
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/repositories"
)

// ListCharactersUseCase handles listing all characters
type ListCharactersUseCase struct {
	characterRepo repositories.CharacterRepository
}

func NewListCharactersUseCase(characterRepo repositories.CharacterRepository) *ListCharactersUseCase {
	return &ListCharactersUseCase{
		characterRepo: characterRepo,
	}
}

// Execute returns a list of all character names
func (uc *ListCharactersUseCase) Execute() (*dtos.CharacterListDTO, error) {
	names, err := uc.characterRepo.List()
	if err != nil {
		return nil, err
	}
	return dtos.ToCharacterListDTO(names), nil
}
