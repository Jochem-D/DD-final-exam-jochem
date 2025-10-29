package usecases

import (
	"strings"

	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/services"
	"ddsheetfinal/internal/domain/valueobjects"
)

// GetEnrichedCharacterUseCase returns a character with all calculated stats
type GetEnrichedCharacterUseCase struct {
	characterRepo    repositories.CharacterRepository
	characterService *services.CharacterService
}

// NewGetEnrichedCharacterUseCase creates a new use case
func NewGetEnrichedCharacterUseCase(
	characterRepo repositories.CharacterRepository,
	characterService *services.CharacterService,
) *GetEnrichedCharacterUseCase {
	return &GetEnrichedCharacterUseCase{
		characterRepo:    characterRepo,
		characterService: characterService,
	}
}

// Execute retrieves and enriches a character with all calculated data
func (uc *GetEnrichedCharacterUseCase) Execute(characterName string) (*dtos.EnrichedCharacterDTO, error) {
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

	// Calculate proficiency bonus
	profBonus := valueobjects.ProficiencyBonus(character.Level)

	// Calculate AC
	ac, acCalc := uc.characterService.ComputeArmorClass(character)

	// Calculate initiative
	initiative := uc.characterService.ComputeInitiativeBonus(character)

	// Calculate passive perception
	passivePerception := uc.characterService.ComputePassivePerception(character)

	// Determine speed (racial speed - simplified)
	speed := uc.getSpeed(character.Race)

	// Calculate saving throws with proficiency
	strSaveProf := uc.characterService.IsSaveProficient(character, "Strength")
	dexSaveProf := uc.characterService.IsSaveProficient(character, "Dexterity")
	conSaveProf := uc.characterService.IsSaveProficient(character, "Constitution")
	intSaveProf := uc.characterService.IsSaveProficient(character, "Intelligence")
	wisSaveProf := uc.characterService.IsSaveProficient(character, "Wisdom")
	chaSaveProf := uc.characterService.IsSaveProficient(character, "Charisma")

	strSave := strMod
	if strSaveProf {
		strSave += profBonus
	}
	dexSave := dexMod
	if dexSaveProf {
		dexSave += profBonus
	}
	conSave := conMod
	if conSaveProf {
		conSave += profBonus
	}
	intSave := intMod
	if intSaveProf {
		intSave += profBonus
	}
	wisSave := wisMod
	if wisSaveProf {
		wisSave += profBonus
	}
	chaSave := chaMod
	if chaSaveProf {
		chaSave += profBonus
	}

	// Calculate skills with proficiency
	skills, skillProfs := uc.calculateSkills(character, strMod, dexMod, conMod, intMod, wisMod, chaMod, profBonus)

	// Build enriched DTO
	enriched := &dtos.EnrichedCharacterDTO{
		CharacterDTO:      *dtos.ToCharacterDTO(character),
		StrMod:            strMod,
		DexMod:            dexMod,
		ConMod:            conMod,
		IntMod:            intMod,
		WisMod:            wisMod,
		ChaMod:            chaMod,
		ProficiencyBonus:  profBonus,
		ArmorClass:        ac,
		ACCalculation:     acCalc,
		Initiative:        initiative,
		PassivePerception: passivePerception,
		Speed:             speed,
		StrSave:           strSave,
		DexSave:           dexSave,
		ConSave:           conSave,
		IntSave:           intSave,
		WisSave:           wisSave,
		ChaSave:           chaSave,
		StrSaveProf:       strSaveProf,
		DexSaveProf:       dexSaveProf,
		ConSaveProf:       conSaveProf,
		IntSaveProf:       intSaveProf,
		WisSaveProf:       wisSaveProf,
		ChaSaveProf:       chaSaveProf,
		Skills:            skills,
		SkillProfs:        skillProfs,
	}

	return enriched, nil
}

func (uc *GetEnrichedCharacterUseCase) getSpeed(race string) int {
	raceLower := strings.ToLower(strings.TrimSpace(race))
	
	// Racial speeds from D&D 5e
	switch raceLower {
	case "dwarf", "hill dwarf", "mountain dwarf":
		return 25
	case "halfling", "lightfoot halfling", "stout halfling":
		return 25
	case "gnome", "forest gnome", "rock gnome":
		return 25
	default:
		return 30 // Default for most races
	}
}

func (uc *GetEnrichedCharacterUseCase) calculateSkills(
	character *entities.Character,
	strMod, dexMod, conMod, intMod, wisMod, chaMod int,
	profBonus int,
) (map[string]int, map[string]bool) {
	// Skill to ability mapping
	skillAbilities := map[string]int{
		"acrobatics":      dexMod,
		"animal handling": wisMod,
		"arcana":          intMod,
		"athletics":       strMod,
		"deception":       chaMod,
		"history":         intMod,
		"insight":         wisMod,
		"intimidation":    chaMod,
		"investigation":   intMod,
		"medicine":        wisMod,
		"nature":          intMod,
		"perception":      wisMod,
		"performance":     chaMod,
		"persuasion":      chaMod,
		"religion":        intMod,
		"sleight of hand": dexMod,
		"stealth":         dexMod,
		"survival":        wisMod,
	}

	skills := make(map[string]int)
	skillProfs := make(map[string]bool)

	// Calculate modifier for each skill
	for skillName, abilityMod := range skillAbilities {
		isProficient := character.HasSkillProficiency(skillName)
		skillProfs[skillName] = isProficient
		
		modifier := abilityMod
		if isProficient {
			modifier += profBonus
		}
		skills[skillName] = modifier
	}

	return skills, skillProfs
}
