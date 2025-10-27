package cli

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/application/usecases"
)

// UnequipHandler handles the unequip command
type UnequipHandler struct {
	unequipItemUseCase *usecases.UnequipItemUseCase
}

// NewUnequipHandler creates a new unequip handler
func NewUnequipHandler(unequipItemUseCase *usecases.UnequipItemUseCase) *UnequipHandler {
	return &UnequipHandler{
		unequipItemUseCase: unequipItemUseCase,
	}
}

// Handle processes the unequip command
func (h *UnequipHandler) Handle(args []string) {
	// Parse input from args
	input := h.parseInput(args)
	if input.CharacterName == "" {
		fmt.Println("name is required")
		os.Exit(2)
	}

	// Execute use case
	if err := h.unequipItemUseCase.Execute(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Equipment updated successfully!")
}

func (h *UnequipHandler) parseInput(args []string) usecases.UnequipItemInput {
	input := usecases.UnequipItemInput{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-name":
			if i+1 < len(args) {
				input.CharacterName = args[i+1]
				i++
			}
		case "-weapon":
			input.UnequipWeapon = true
		case "-armor":
			input.UnequipArmor = true
		case "-shield":
			input.UnequipShield = true
		}
	}

	return input
}
