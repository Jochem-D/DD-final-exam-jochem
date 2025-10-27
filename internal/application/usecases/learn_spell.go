package usecases

import (
	"fmt"
	"sort"
	"strings"

	"ddsheetfinal/internal/domain/repositories"
)

// LearnSpellUseCase handles learning new spells
type LearnSpellUseCase struct {
	characterRepo repositories.CharacterRepository
	srdRepo       repositories.SRDRepository
}

// NewLearnSpellUseCase creates a new LearnSpellUseCase
func NewLearnSpellUseCase(
	characterRepo repositories.CharacterRepository,
	srdRepo repositories.SRDRepository,
) *LearnSpellUseCase {
	return &LearnSpellUseCase{
		characterRepo: characterRepo,
		srdRepo:       srdRepo,
	}
}

// Execute learns a new spell
func (uc *LearnSpellUseCase) Execute(characterName, spellName string) error {
	// Load character
	character, err := uc.characterRepo.FindByName(characterName)
	if err != nil {
		return err
	}

	classLower := strings.ToLower(strings.TrimSpace(character.Class))

	// Check if class can cast spells
	if !uc.isCasterClass(classLower) {
		return fmt.Errorf("this class can't cast spells")
	}

	// Check if class learns spells (vs prepares them)
	if uc.isPreparedClass(classLower) {
		return fmt.Errorf("this class prepares spells and can't learn them")
	}

	// Validate spell exists
	if !uc.srdRepo.IsValidSpell(spellName) {
		return fmt.Errorf("spell '%s' not found in 5e-SRD-Spells.csv", spellName)
	}

	// Validate spell is available to this class
	ok, err := uc.srdRepo.IsSpellForClass(spellName, character.Class)
	if err != nil {
		return fmt.Errorf("error validating spell for class: %w", err)
	}
	if !ok {
		return fmt.Errorf("'%s' is not available to the %s class", spellName, classLower)
	}

	// Check if already learned
	for _, s := range character.Spells {
		if strings.EqualFold(s, spellName) {
			return fmt.Errorf("%s already knows the spell '%s'", character.Name, spellName)
		}
	}

	// Learn the spell
	character.Spells = append(character.Spells, spellName)
	sort.Strings(character.Spells)

	// Save character
	return uc.characterRepo.Save(character)
}

func (uc *LearnSpellUseCase) isCasterClass(class string) bool {
	casters := []string{"bard", "cleric", "druid", "paladin", "ranger", "sorcerer", "warlock", "wizard"}
	for _, c := range casters {
		if c == class {
			return true
		}
	}
	return false
}

func (uc *LearnSpellUseCase) isPreparedClass(class string) bool {
	prepared := []string{"cleric", "druid", "paladin", "wizard"}
	for _, c := range prepared {
		if c == class {
			return true
		}
	}
	return false
}
