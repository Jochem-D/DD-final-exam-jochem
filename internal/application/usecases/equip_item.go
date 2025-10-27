package usecases

import (
	"fmt"
	"strings"

	"ddsheetfinal/internal/domain/repositories"
)

// EquipItemInput contains the input for equipping items
type EquipItemInput struct {
	CharacterName string
	Weapon        string
	Armor         string
	Shield        string
	Slot          string // "main hand" or "off hand" for weapons
}

// EquipItemUseCase handles equipping items to a character
type EquipItemUseCase struct {
	characterRepo repositories.CharacterRepository
	srdRepo       repositories.SRDRepository
}

// NewEquipItemUseCase creates a new EquipItemUseCase
func NewEquipItemUseCase(
	characterRepo repositories.CharacterRepository,
	srdRepo repositories.SRDRepository,
) *EquipItemUseCase {
	return &EquipItemUseCase{
		characterRepo: characterRepo,
		srdRepo:       srdRepo,
	}
}

// Execute equips items to a character
func (uc *EquipItemUseCase) Execute(input EquipItemInput) error {
	// Load character
	character, err := uc.characterRepo.FindByName(input.CharacterName)
	if err != nil {
		return err
	}

	changed := false

	// Equip weapon
	if input.Weapon != "" {
		if !uc.srdRepo.IsValidEquipment(input.Weapon) {
			return fmt.Errorf("weapon '%s' not found in 5e-SRD-Equipment.csv", input.Weapon)
		}

		slot := uc.normalizeSlot(input.Slot)
		if slot == "" {
			slot = "main hand"
		}

		if slot == "main hand" {
			if character.Weapon != "" {
				return fmt.Errorf("main hand already occupied")
			}
			character.Weapon = input.Weapon
		} else if slot == "off hand" {
			if character.OffHand != "" {
				return fmt.Errorf("off hand already occupied")
			}
			character.OffHand = input.Weapon
		} else {
			return fmt.Errorf("invalid slot: %s", input.Slot)
		}
		changed = true
	}

	// Equip armor
	if input.Armor != "" {
		if !uc.srdRepo.IsValidEquipment(input.Armor) {
			return fmt.Errorf("armor '%s' not found in 5e-SRD-Equipment.csv", input.Armor)
		}
		if character.Armor != "" {
			return fmt.Errorf("armor slot already occupied")
		}
		character.Armor = input.Armor
		changed = true
	}

	// Equip shield
	if input.Shield != "" {
		if !uc.srdRepo.IsValidEquipment(input.Shield) {
			return fmt.Errorf("shield '%s' not found in 5e-SRD-Equipment.csv", input.Shield)
		}
		if character.Shield != "" {
			return fmt.Errorf("shield slot already occupied")
		}
		character.Shield = input.Shield
		changed = true
	}

	if !changed {
		return fmt.Errorf("no equipment was specified")
	}

	// Save character
	return uc.characterRepo.Save(character)
}

func (uc *EquipItemUseCase) normalizeSlot(slot string) string {
	s := strings.ToLower(strings.TrimSpace(slot))
	switch s {
	case "main", "main hand", "main-hand", "mainhand":
		return "main hand"
	case "off", "off hand", "off-hand", "offhand":
		return "off hand"
	default:
		return s
	}
}
