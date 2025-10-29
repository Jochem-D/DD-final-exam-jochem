package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	enrichmentRepo   repositories.EnrichmentRepository
}

// NewGetEnrichedCharacterUseCase creates a new use case
func NewGetEnrichedCharacterUseCase(
	characterRepo repositories.CharacterRepository,
	characterService *services.CharacterService,
	enrichmentRepo repositories.EnrichmentRepository,
) *GetEnrichedCharacterUseCase {
	return &GetEnrichedCharacterUseCase{
		characterRepo:    characterRepo,
		characterService: characterService,
		enrichmentRepo:   enrichmentRepo,
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

	// Calculate weapon attacks with enrichment data
	weaponAttacks := uc.calculateWeaponAttacks(character, strMod, dexMod, profBonus)

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
		WeaponAttacks:     weaponAttacks,
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

func (uc *GetEnrichedCharacterUseCase) calculateWeaponAttacks(
	character *entities.Character,
	strMod, dexMod int,
	profBonus int,
) []dtos.WeaponAttackDTO {
	attacks := []dtos.WeaponAttackDTO{}

	// Add weapon if equipped
	if character.Weapon != "" {
		attack := uc.calculateSingleWeaponAttack(character.Weapon, strMod, dexMod, profBonus, character.Name)
		if attack != nil {
			attacks = append(attacks, *attack)
		}
	}

	// Add off-hand if equipped
	if character.OffHand != "" {
		attack := uc.calculateSingleWeaponAttack(character.OffHand, strMod, dexMod, profBonus, character.Name)
		if attack != nil {
			attacks = append(attacks, *attack)
		}
	}

	return attacks
}

func (uc *GetEnrichedCharacterUseCase) calculateSingleWeaponAttack(
	weaponName string,
	strMod, dexMod int,
	profBonus int,
	characterName string,
) *dtos.WeaponAttackDTO {
	if weaponName == "" {
		return nil
	}

	// Try to load enriched data from JSON file
	enrichedPath := filepath.Join("data", "enrichments", characterName+".json")
	enrichedData := make(map[string]interface{})
	
	if data, err := os.ReadFile(enrichedPath); err == nil {
		json.Unmarshal(data, &enrichedData)
	}

	// Try to get weapon info from enriched.equipment
	var damageText, properties string
	if enriched, ok := enrichedData["enriched"].(map[string]interface{}); ok {
		if equipment, ok := enriched["equipment"].(map[string]interface{}); ok {
			if weaponData, ok := equipment[weaponName].(map[string]interface{}); ok {
				if dt, ok := weaponData["damage_text"].(string); ok {
					damageText = dt
				}
				if props, ok := weaponData["properties"].(string); ok {
					properties = props
				}
			}
		}
	}

	// If no enriched data, use basic weapon data as fallback
	if damageText == "" {
		return uc.calculateBasicWeaponAttack(weaponName, strMod, dexMod, profBonus)
	}

	// Determine which ability mod to use
	abilityMod := strMod
	isFinesse := strings.Contains(strings.ToLower(properties), "finesse")
	if isFinesse && dexMod > strMod {
		abilityMod = dexMod
	}

	// Attack bonus = ability mod + proficiency bonus
	attackBonus := abilityMod + profBonus

	// Format attack bonus with sign
	attackBonusStr := formatSigned(attackBonus)

	// Format damage: dice + ability mod (e.g., "1d6+5")
	damage := damageText
	if abilityMod != 0 {
		damage = damage + formatSigned(abilityMod)
	}

	return &dtos.WeaponAttackDTO{
		Name:        weaponName,
		AttackBonus: attackBonusStr,
		Damage:      damage,
	}
}

func (uc *GetEnrichedCharacterUseCase) calculateBasicWeaponAttack(
	weaponName string,
	strMod, dexMod int,
	profBonus int,
) *dtos.WeaponAttackDTO {
	// Fallback when enrichment data is not available
	// Use STR mod by default for unknown weapons
	abilityMod := strMod
	attackBonus := abilityMod + profBonus

	return &dtos.WeaponAttackDTO{
		Name:        weaponName,
		AttackBonus: formatSigned(attackBonus),
		Damage:      "unknown (enrich character first)",
	}
}

func formatSigned(n int) string {
	if n >= 0 {
		return fmt.Sprintf("+%d", n)
	}
	return fmt.Sprintf("%d", n)
}
