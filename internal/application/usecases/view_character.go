package usecases

import (
	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/services"
)

// ViewCharacterOutput contains character data and computed derived stats
type ViewCharacterOutput struct {
	Character        *entities.Character
	ArmorClass       int
	ArmorClassDesc   string
	InitiativeBonus  int
	PassivePerception int
}

// ViewCharacterUseCase handles viewing character details
type ViewCharacterUseCase struct {
	characterRepo  repositories.CharacterRepository
	characterService *services.CharacterService
}

// NewViewCharacterUseCase creates a new ViewCharacterUseCase
func NewViewCharacterUseCase(
	characterRepo repositories.CharacterRepository,
	characterService *services.CharacterService,
) *ViewCharacterUseCase {
	return &ViewCharacterUseCase{
		characterRepo:  characterRepo,
		characterService: characterService,
	}
}

// Execute retrieves a character and calculates derived stats
func (uc *ViewCharacterUseCase) Execute(name string) (*ViewCharacterOutput, error) {
	character, err := uc.characterRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	// Calculate derived stats
	ac, acDesc := uc.characterService.ComputeArmorClass(character)
	initiative := uc.characterService.ComputeInitiativeBonus(character)
	passivePerception := uc.characterService.ComputePassivePerception(character)

	return &ViewCharacterOutput{
		Character:        character,
		ArmorClass:       ac,
		ArmorClassDesc:   acDesc,
		InitiativeBonus:  initiative,
		PassivePerception: passivePerception,
	}, nil
}
