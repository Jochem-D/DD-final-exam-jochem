package usecases

import (
	"fmt"
	"sort"
	"strings"

	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
	"ddsheetfinal/internal/domain/valueobjects"
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
	srdRepo       repositories.SRDRepository
}

func NewPrepareSpellUseCase(characterRepo repositories.CharacterRepository, srdRepo repositories.SRDRepository) *PrepareSpellUseCase {
	return &PrepareSpellUseCase{
		characterRepo: characterRepo,
		srdRepo:       srdRepo,
	}
}

// Execute prepares or unprepares a spell
func (uc *PrepareSpellUseCase) Execute(input PrepareSpellInput) error {
	// Load character
	character, err := uc.characterRepo.FindByName(input.CharacterName)
	if err != nil {
		return err
	}

	// Check if class can cast spells
	classLower := strings.ToLower(character.Class)
	if !valueobjects.IsSpellcaster(classLower) {
		return fmt.Errorf("this class can't cast spells")
	}

	// Check if class prepares spells (vs learning them)
	if !valueobjects.IsPreparedCaster(classLower) {
		return fmt.Errorf("this class learns spells and can't prepare them")
	}

	// For prepared casters, check if spell is valid for their class
	isValid, err := uc.srdRepo.IsSpellForClass(input.SpellName, character.Class)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("spell '%s' is not available to the %s class", input.SpellName, character.Class)
	}

	// Check if character has spell slots for this level
	spellLevel, err := uc.srdRepo.GetSpellLevel(input.SpellName)
	if err != nil {
		return fmt.Errorf("error getting spell level: %w", err)
	}
	
	// Get character's available spell slots
	spellSlots := valueobjects.GetSpellSlots(classLower, character.Level)
	if spellSlots == nil {
		return fmt.Errorf("this class can't cast spells")
	}
	
	// Check if they have slots for this spell level (skip cantrips at index 0)
	if spellLevel > 0 && spellLevel < len(spellSlots) {
		if spellSlots[spellLevel] == 0 {
			return fmt.Errorf("the spell has higher level than the available spell slots")
		}
	} else if spellLevel >= len(spellSlots) {
		return fmt.Errorf("the spell has higher level than the available spell slots")
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
