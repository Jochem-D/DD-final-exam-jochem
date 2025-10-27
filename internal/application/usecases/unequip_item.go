package usecases

import (
	"fmt"

	"ddsheetfinal/internal/domain/repositories"
)

// UnequipItemInput contains the input for unequipping items
type UnequipItemInput struct {
	CharacterName string
	UnequipWeapon bool
	UnequipArmor  bool
	UnequipShield bool
}

// UnequipItemUseCase handles unequipping items from a character
type UnequipItemUseCase struct {
	characterRepo repositories.CharacterRepository
}

// NewUnequipItemUseCase creates a new UnequipItemUseCase
func NewUnequipItemUseCase(characterRepo repositories.CharacterRepository) *UnequipItemUseCase {
	return &UnequipItemUseCase{
		characterRepo: characterRepo,
	}
}

// Execute unequips items from a character
func (uc *UnequipItemUseCase) Execute(input UnequipItemInput) error {
	// Load character
	character, err := uc.characterRepo.FindByName(input.CharacterName)
	if err != nil {
		return err
	}

	changed := false

	// Unequip weapon
	if input.UnequipWeapon && character.Weapon != "" {
		character.Inventory = append(character.Inventory, character.Weapon)
		character.Weapon = ""
		changed = true
	}

	// Unequip armor
	if input.UnequipArmor && character.Armor != "" {
		character.Inventory = append(character.Inventory, character.Armor)
		character.Armor = ""
		changed = true
	}

	// Unequip shield
	if input.UnequipShield && character.Shield != "" {
		character.Inventory = append(character.Inventory, character.Shield)
		character.Shield = ""
		changed = true
	}

	if !changed {
		return fmt.Errorf("no equipment was unequipped")
	}

	// Save character
	return uc.characterRepo.Save(character)
}
