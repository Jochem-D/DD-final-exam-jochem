package usecases

import (
	"strings"
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/services"
	"ddsheetfinal/internal/domain/valueobjects"
)

// ViewCharacterOutput contains character data and computed derived stats
type ViewCharacterOutput struct {
	Character         *dtos.CharacterDTO
	ArmorClass        int
	ArmorClassDesc    string
	InitiativeBonus   int
	PassivePerception int
	
	// Ability modifiers
	StrMod int
	DexMod int
	ConMod int
	IntMod int
	WisMod int
	ChaMod int
	
	// Proficiency bonus
	ProficiencyBonus int
	
	// Spell information (nil if not a caster)
	SpellSlots         []int  // Slots per level
	SpellcastingAbility string // e.g., "intelligence", "wisdom", "charisma"
	SpellSaveDC        int
	SpellAttackBonus   int
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
	
	// Calculate ability modifiers
	strMod := valueobjects.AbilityModifier(character.Str)
	dexMod := valueobjects.AbilityModifier(character.Dex)
	conMod := valueobjects.AbilityModifier(character.Con)
	intMod := valueobjects.AbilityModifier(character.Int)
	wisMod := valueobjects.AbilityModifier(character.Wis)
	chaMod := valueobjects.AbilityModifier(character.Cha)
	
	// Calculate proficiency bonus
	profBonus := valueobjects.ProficiencyBonus(character.Level)
	
	// Calculate spell information
	classLower := strings.ToLower(character.Class)
	spellSlots := valueobjects.GetSpellSlots(classLower, character.Level)
	var spellcastingAbility string
	var spellSaveDC, spellAttackBonus int
	
	if spellSlots != nil {
		abilityAbbr := valueobjects.GetSpellcastingAbility(classLower)
		if abilityAbbr != "" {
			abilityScore := character.GetAbilityScore(abilityAbbr)
			abilityMod := valueobjects.AbilityModifier(abilityScore)
			spellSaveDC = 8 + profBonus + abilityMod
			spellAttackBonus = profBonus + abilityMod
			
			// Convert abbreviation to full name
			switch abilityAbbr {
			case "int":
				spellcastingAbility = "intelligence"
			case "wis":
				spellcastingAbility = "wisdom"
			case "cha":
				spellcastingAbility = "charisma"
			}
		}
	}

	return &ViewCharacterOutput{
		Character:          dtos.ToCharacterDTO(character),
		ArmorClass:         ac,
		ArmorClassDesc:     acDesc,
		InitiativeBonus:    initiative,
		PassivePerception:  passivePerception,
		StrMod:             strMod,
		DexMod:             dexMod,
		ConMod:             conMod,
		IntMod:             intMod,
		WisMod:             wisMod,
		ChaMod:             chaMod,
		ProficiencyBonus:   profBonus,
		SpellSlots:         spellSlots,
		SpellcastingAbility: spellcastingAbility,
		SpellSaveDC:        spellSaveDC,
		SpellAttackBonus:   spellAttackBonus,
	}, nil
}
