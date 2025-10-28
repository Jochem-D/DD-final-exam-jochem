package usecases

import (
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/services"
	"ddsheetfinal/internal/domain/valueobjects"
)

// DeriveCharacterStatsUseCase calculates derived stats for a character
type DeriveCharacterStatsUseCase struct {
	characterRepo    repositories.CharacterRepository
	characterService *services.CharacterService
}

// NewDeriveCharacterStatsUseCase creates a new DeriveCharacterStatsUseCase
func NewDeriveCharacterStatsUseCase(
	characterRepo repositories.CharacterRepository,
	characterService *services.CharacterService,
) *DeriveCharacterStatsUseCase {
	return &DeriveCharacterStatsUseCase{
		characterRepo:    characterRepo,
		characterService: characterService,
	}
}

// Execute calculates derived stats for a character
func (uc *DeriveCharacterStatsUseCase) Execute(characterName string) (*dtos.DerivedStatsDTO, error) {
	// Load character
	character, err := uc.characterRepo.FindByName(characterName)
	if err != nil {
		return nil, err
	}

	// Calculate ability modifiers
	strMod := valueobjects.AbilityModifier(character.Str)
	dexMod := valueobjects.AbilityModifier(character.Dex)
	conMod := valueobjects.AbilityModifier(character.Con)
	intMod := valueobjects.AbilityModifier(character.Int)
	wisMod := valueobjects.AbilityModifier(character.Wis)
	chaMod := valueobjects.AbilityModifier(character.Cha)

	// Calculate AC
	ac, acCalc := uc.characterService.ComputeArmorClass(character)

	// Calculate proficiency bonus
	profBonus := valueobjects.ProficiencyBonus(character.Level)

	// Calculate saving throws
	strSave := strMod
	dexSave := dexMod
	conSave := conMod
	intSave := intMod
	wisSave := wisMod
	chaSave := chaMod
	
	if uc.characterService.IsSaveProficient(character, "Strength") {
		strSave += profBonus
	}
	if uc.characterService.IsSaveProficient(character, "Dexterity") {
		dexSave += profBonus
	}
	if uc.characterService.IsSaveProficient(character, "Constitution") {
		conSave += profBonus
	}
	if uc.characterService.IsSaveProficient(character, "Intelligence") {
		intSave += profBonus
	}
	if uc.characterService.IsSaveProficient(character, "Wisdom") {
		wisSave += profBonus
	}
	if uc.characterService.IsSaveProficient(character, "Charisma") {
		chaSave += profBonus
	}
	
	saveStr := "None" // Simplified

	// Calculate skills - simplified, just return empty map for now
	skills := make(map[string]int)

	return &dtos.DerivedStatsDTO{
		CharacterName:    characterName,
		StrMod:           strMod,
		DexMod:           dexMod,
		ConMod:           conMod,
		IntMod:           intMod,
		WisMod:           wisMod,
		ChaMod:           chaMod,
		AC:               ac,
		ACCalculation:    acCalc,
		StrSave:          strSave,
		DexSave:          dexSave,
		ConSave:          conSave,
		IntSave:          intSave,
		WisSave:          wisSave,
		ChaSave:          chaSave,
		SavingThrowStr:   saveStr,
		Skills:           skills,
		ProficiencyBonus: profBonus,
	}, nil
}
