package usecases

import (
	"ddsheetfinal/internal/domain/repositories"
)

// GetLearnableSpellsUseCase handles retrieving spells a character can learn
type GetLearnableSpellsUseCase struct {
	characterRepo repositories.CharacterRepository
	srdRepo       repositories.SRDRepository
}

// NewGetLearnableSpellsUseCase creates a new GetLearnableSpellsUseCase
func NewGetLearnableSpellsUseCase(
	characterRepo repositories.CharacterRepository,
	srdRepo repositories.SRDRepository,
) *GetLearnableSpellsUseCase {
	return &GetLearnableSpellsUseCase{
		characterRepo: characterRepo,
		srdRepo:       srdRepo,
	}
}

// Execute returns spells the character can still learn
func (uc *GetLearnableSpellsUseCase) Execute(characterName string) ([]string, error) {
	// Load character
	character, err := uc.characterRepo.FindByName(characterName)
	if err != nil {
		return nil, err
	}

	// Get learnable spells
	return uc.srdRepo.GetLearnableSpells(character.Class, character.Spells)
}
