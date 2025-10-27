package usecases

import (
	"fmt"
	"sort"
	"strings"

	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
)

// PrepareSpellInput contains input for preparing/unpreparing spells
type PrepareSpellInput struct {
	CharacterName string
	SpellName     string
	Remove        bool // true to unprepare, false to prepare
}

// PrepareSpellUseCase handles preparing and unpreparing spells
type PrepareSpellUseCase struct {
	characterRepo repositories.CharacterRepository
}

func NewPrepareSpellUseCase(characterRepo repositories.CharacterRepository) *PrepareSpellUseCase {
	return &PrepareSpellUseCase{
		characterRepo: characterRepo,
	}
}

// Execute prepares or unprepares a spell
func (uc *PrepareSpellUseCase) Execute(input PrepareSpellInput) error {
	// Load character
	character, err := uc.characterRepo.FindByName(input.CharacterName)
	if err != nil {
		return err
	}

	// Check if spell is in the character's learned spells
	spellLearned := false
	for _, s := range character.Spells {
		if strings.EqualFold(s, input.SpellName) {
			spellLearned = true
			break
		}
	}

	if !spellLearned {
		return fmt.Errorf("spell '%s' is not in %s's learned spells", input.SpellName, character.Name)
	}

	if input.Remove {
		// Unprepare the spell
		return uc.unprepareSpell(character, input.SpellName)
	}

	// Prepare the spell
	return uc.prepareSpell(character, input.SpellName)
}

func (uc *PrepareSpellUseCase) prepareSpell(character *entities.Character, spellName string) error {
	// Check if already prepared
	for _, s := range character.PreparedSpells {
		if strings.EqualFold(s, spellName) {
			return fmt.Errorf("spell '%s' is already prepared", spellName)
		}
	}

	// Prepare the spell
	character.PreparedSpells = append(character.PreparedSpells, spellName)
	sort.Strings(character.PreparedSpells)

	return uc.characterRepo.Save(character)
}

func (uc *PrepareSpellUseCase) unprepareSpell(character *entities.Character, spellName string) error {
	// Find and remove the spell
	found := false
	newPrepared := make([]string, 0, len(character.PreparedSpells))
	for _, s := range character.PreparedSpells {
		if !strings.EqualFold(s, spellName) {
			newPrepared = append(newPrepared, s)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("spell '%s' is not prepared", spellName)
	}

	character.PreparedSpells = newPrepared
	return uc.characterRepo.Save(character)
}
