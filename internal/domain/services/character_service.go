package services

import (
	"fmt"
	"strings"

	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/valueobjects"
)

// CharacterService provides domain logic for character operations
type CharacterService struct{}

// NewCharacterService creates a new CharacterService
func NewCharacterService() *CharacterService {
	return &CharacterService{}
}

// ComputeArmorClass calculates the character's AC based on armor, dex, and class
func (s *CharacterService) ComputeArmorClass(char *entities.Character) (int, string) {
	dexMod := valueobjects.AbilityModifier(char.Dex)
	conMod := valueobjects.AbilityModifier(char.Con)
	wisMod := valueobjects.AbilityModifier(char.Wis)
	
	shieldBonus := 0
	if strings.TrimSpace(strings.ToLower(char.Shield)) != "" {
		shieldBonus = 2
	}
	
	armorName := strings.ToLower(strings.TrimSpace(char.Armor))
	
	// If no armor, check for class-specific unarmored defense
	if armorName == "" {
		return s.computeUnarmoredAC(char, dexMod, conMod, wisMod, shieldBonus)
	}
	
	// Look up armor info
	armorInfo, ok := valueobjects.GetArmorInfo(armorName)
	if !ok {
		// Invalid armor, treat as unarmored
		return s.computeUnarmoredAC(char, dexMod, conMod, wisMod, shieldBonus)
	}
	
	return s.computeArmoredAC(armorInfo, dexMod, shieldBonus)
}

func (s *CharacterService) computeUnarmoredAC(char *entities.Character, dexMod, conMod, wisMod, shieldBonus int) (int, string) {
	classLower := strings.ToLower(strings.TrimSpace(char.Class))
	
	switch classLower {
	case "monk":
		if shieldBonus == 0 {
			ac := 10 + dexMod + wisMod
			return ac, fmt.Sprintf("10 + Dex(%+d) + Wis(%+d)", dexMod, wisMod)
		}
	case "barbarian":
		ac := 10 + dexMod + conMod + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("10 + Dex(%+d) + Con(%+d) + shield(+2)", dexMod, conMod)
		}
		return ac, fmt.Sprintf("10 + Dex(%+d) + Con(%+d)", dexMod, conMod)
	}
	
	ac := 10 + dexMod + shieldBonus
	if shieldBonus > 0 {
		return ac, fmt.Sprintf("10 + Dex(%+d) + shield(+2)", dexMod)
	}
	return ac, fmt.Sprintf("10 + Dex(%+d)", dexMod)
}

func (s *CharacterService) computeArmoredAC(armor valueobjects.ArmorInfo, dexMod, shieldBonus int) (int, string) {
	switch armor.Type {
	case valueobjects.ArmorTypeLight:
		ac := armor.Base + dexMod + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("%d (light base) + Dex(%+d) + shield(+2)", armor.Base, dexMod)
		}
		return ac, fmt.Sprintf("%d (light base) + Dex(%+d)", armor.Base, dexMod)
		
	case valueobjects.ArmorTypeMedium:
		cappedDex := dexMod
		if cappedDex > 2 {
			cappedDex = 2
		}
		ac := armor.Base + cappedDex + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("%d (medium base) + Dex cap(+%d) + shield(+2)", armor.Base, cappedDex)
		}
		return ac, fmt.Sprintf("%d (medium base) + Dex cap(+%d)", armor.Base, cappedDex)
		
	case valueobjects.ArmorTypeHeavy:
		ac := armor.Base + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("%d (heavy base) + shield(+2)", armor.Base)
		}
		return ac, fmt.Sprintf("%d (heavy base)", armor.Base)
		
	default:
		ac := 10 + dexMod + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("10 + Dex(%+d) + shield(+2)", dexMod)
		}
		return ac, fmt.Sprintf("10 + Dex(%+d)", dexMod)
	}
}

// ComputeInitiativeBonus calculates the initiative bonus (dex modifier)
func (s *CharacterService) ComputeInitiativeBonus(char *entities.Character) int {
	return valueobjects.AbilityModifier(char.Dex)
}

// ComputePassivePerception calculates passive perception (10 + Wis mod + proficiency if proficient)
func (s *CharacterService) ComputePassivePerception(char *entities.Character) int {
	wisMod := valueobjects.AbilityModifier(char.Wis)
	base := 10 + wisMod
	
	if char.HasSkillProficiency("perception") {
		profBonus := valueobjects.ProficiencyBonus(char.Level)
		base += profBonus
	}
	
	return base
}

// IsSaveProficient checks if the character is proficient in a saving throw
func (s *CharacterService) IsSaveProficient(char *entities.Character, ability string) bool {
	classLower := strings.ToLower(strings.TrimSpace(char.Class))
	ability = strings.ToLower(strings.TrimSpace(ability))
	
	// Normalize ability name to abbreviation
	ability = s.normalizeAbility(ability)
	
	profs, ok := valueobjects.ClassSavingThrowProficiencies[classLower]
	if !ok {
		return false
	}
	
	for _, p := range profs {
		if p == ability {
			return true
		}
	}
	return false
}

func (s *CharacterService) normalizeAbility(ability string) string {
	switch ability {
	case "strength", "str":
		return "str"
	case "dexterity", "dex":
		return "dex"
	case "constitution", "con":
		return "con"
	case "intelligence", "int":
		return "int"
	case "wisdom", "wis":
		return "wis"
	case "charisma", "cha":
		return "cha"
	default:
		return ability
	}
}

// ComputeSkillModifier calculates the modifier for a given skill
func (s *CharacterService) ComputeSkillModifier(char *entities.Character, skill string) int {
	skillLower := strings.ToLower(strings.TrimSpace(skill))
	
	// Get the governing ability for this skill
	ability, ok := valueobjects.SkillToAbility[skillLower]
	if !ok {
		return 0
	}
	
	abilityScore := char.GetAbilityScore(ability)
	abilityMod := valueobjects.AbilityModifier(abilityScore)
	
	// Add proficiency bonus if proficient
	if char.HasSkillProficiency(skillLower) {
		profBonus := valueobjects.ProficiencyBonus(char.Level)
		return abilityMod + profBonus
	}
	
	return abilityMod
}

// ComputeMaxHP calculates maximum hit points based on class, level, and CON (exam rules)
func (s *CharacterService) ComputeMaxHP(char *entities.Character) int {
	if char.Level <= 0 {
		return 0
	}

	conMod := valueobjects.AbilityModifier(char.Con)

	// After level 1: fixed average
	fixed := s.getHitDieForClass(char.Class) // 7/6/5/4
	// Level 1: max die value per class
	maxDie := s.getMaxDieForClass(char.Class) // 12/10/8/6

	// Level 1
	hp := maxDie + conMod

	// Levels 2..N
	if char.Level > 1 {
		hp += (char.Level - 1) * (fixed + conMod)
	}

	if hp < 1 {
		hp = 1
	}
	return hp
}


// getHitDieForClass returns the fixed average HP per level for each class
func (s *CharacterService) getHitDieForClass(class string) int {
	classLower := strings.ToLower(strings.TrimSpace(class))
	
	hitDieMap := map[string]int{
		"barbarian": 7, // d12
		"fighter":   6, // d10
		"paladin":   6, // d10
		"ranger":    6, // d10
		"bard":      5, // d8
		"cleric":    5, // d8
		"druid":     5, // d8
		"monk":      5, // d8
		"rogue":     5, // d8
		"warlock":   5, // d8
		"sorcerer":  4, // d6
		"wizard":    4, // d6
	}
	
	if hitDie, ok := hitDieMap[classLower]; ok {
		return hitDie
	}
	
	// Default to d8 (5) if class is not found
	return 5
}

// getMaxDieForClass returns the max hit die value used at level 1
func (s *CharacterService) getMaxDieForClass(class string) int {
	classLower := strings.ToLower(strings.TrimSpace(class))

	switch classLower {
	case "barbarian":
		return 12
	case "fighter", "paladin", "ranger":
		return 10
	case "wizard", "sorcerer":
		return 6
	default: // d8 bucket (rogue, bard, cleric, druid, monk, warlock, artificer, etc.)
		return 8
	}
}