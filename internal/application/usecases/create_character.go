package usecases

import (
	"fmt"
	"sort"
	"strings"

	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/valueobjects"
)

type CreateCharacterUseCase struct {
	characterRepo repositories.CharacterRepository
}

func NewCreateCharacterUseCase(characterRepo repositories.CharacterRepository) *CreateCharacterUseCase {
	return &CreateCharacterUseCase{
		characterRepo: characterRepo,
	}
}

// Execute creates a new character
func (uc *CreateCharacterUseCase) Execute(input dtos.CreateCharacterDTO) (*dtos.CharacterDTO, error) {
	// Validate input
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.Race == "" {
		return nil, fmt.Errorf("race is required")
	}
	if input.Class == "" {
		return nil, fmt.Errorf("class is required")
	}

	// Build skill list
	skillList, err := uc.buildSkillList(input)
	if err != nil {
		return nil, err
	}

	// Apply racial bonuses
	str, dex, con, intScore, wis, cha := uc.applyRacialBonuses(input)

	// Sort skills
	sort.Strings(skillList)

	// Create character entity
	character := &entities.Character{
		Name:               input.Name,
		Race:               input.Race,
		Class:              input.Class,
		Background:         input.Background,
		Level:              input.Level,
		Str:                str,
		Dex:                dex,
		Con:                con,
		Int:                intScore,
		Wis:                wis,
		Cha:                cha,
		SkillProficiencies: skillList,
	}

	// Save character
	if err := uc.characterRepo.Save(character); err != nil {
		return nil, err
	}

	return dtos.ToCharacterDTO(character), nil
}

func (uc *CreateCharacterUseCase) buildSkillList(input dtos.CreateCharacterDTO) ([]string, error) {
	classLower := strings.ToLower(input.Class)
	backgroundLower := strings.ToLower(input.Background)

	// Start with background skills (allow duplicates for correct output)
	bgSkills := valueobjects.BackgroundSkillProficiencies[backgroundLower]
	result := make([]string, 0)
	for _, s := range bgSkills {
		result = append(result, s)
	}

	// If skills are provided, validate and add them
	if len(input.Skills) > 0 {
		classSkills := valueobjects.ClassSkillProficiencies[classLower]
		classSkillMap := make(map[string]bool)
		for _, s := range classSkills {
			classSkillMap[s] = true
		}

		for _, skill := range input.Skills {
			skillLower := strings.ToLower(strings.TrimSpace(skill))
			if !classSkillMap[skillLower] {
				return nil, fmt.Errorf("skill '%s' is not available to the %s class", skill, input.Class)
			}
			result = append(result, skillLower)
		}

		// Check number of skills
		numClassSkills := valueobjects.ClassSkillChoices[classLower]
		providedSkills := len(input.Skills)
		if providedSkills != numClassSkills {
			return nil, fmt.Errorf("class %s requires exactly %d skill proficiencies, got %d", input.Class, numClassSkills, providedSkills)
		}
	} else {
		// Auto-assign class skills
		classSkills := valueobjects.ClassSkillProficiencies[classLower]
		numNeeded := valueobjects.ClassSkillChoices[classLower]

		assigned := 0
		for _, skill := range classSkills {
			// Don't check for duplicates - add skills even if background already has them
			if assigned < numNeeded {
				result = append(result, skill)
				assigned++
			}
			if assigned >= numNeeded {
				break
			}
		}
	}

	// Sort the list
	sort.Strings(result)

	return result, nil
}

func (uc *CreateCharacterUseCase) applyRacialBonuses(input dtos.CreateCharacterDTO) (str, dex, con, intScore, wis, cha int) {
	str, dex, con, intScore, wis, cha = input.Str, input.Dex, input.Con, input.Int, input.Wis, input.Cha
	raceLower := strings.ToLower(input.Race)

	bonuses, ok := valueobjects.RacialAbilityBonuses[raceLower]
	if !ok {
		return
	}

	for ability, bonus := range bonuses {
		switch ability {
		case "str":
			str += bonus
		case "dex":
			dex += bonus
		case "con":
			con += bonus
		case "int":
			intScore += bonus
		case "wis":
			wis += bonus
		case "cha":
			cha += bonus
		}
	}

	return
}
