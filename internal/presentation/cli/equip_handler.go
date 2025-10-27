package cli

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/application/usecases"
)

// EquipHandler handles the equip command
type EquipHandler struct {
	equipItemUseCase *usecases.EquipItemUseCase
}

// NewEquipHandler creates a new equip handler
func NewEquipHandler(equipItemUseCase *usecases.EquipItemUseCase) *EquipHandler {
	return &EquipHandler{
		equipItemUseCase: equipItemUseCase,
	}
}

// Handle processes the equip command
func (h *EquipHandler) Handle(args []string) {
	// Parse input from args
	input := h.parseInput(args)
	if input.CharacterName == "" {
		fmt.Println("name is required")
		os.Exit(2)
	}

	// Execute use case
	if err := h.equipItemUseCase.Execute(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Print success messages based on what was equipped
	if input.Weapon != "" {
		slot := input.Slot
		if slot == "" {
			slot = "main hand"
		}
		fmt.Printf("Equipped weapon %s to %s\n", input.Weapon, slot)
	}
	if input.Armor != "" {
		fmt.Printf("Equipped armor %s\n", input.Armor)
	}
	if input.Shield != "" {
		fmt.Printf("Equipped shield %s\n", input.Shield)
	}
}

func (h *EquipHandler) parseInput(args []string) usecases.EquipItemInput {
	input := usecases.EquipItemInput{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-name":
			if i+1 < len(args) {
				input.CharacterName = args[i+1]
				i++
			}
		case "-weapon":
			if i+1 < len(args) {
				input.Weapon = args[i+1]
				i++
			}
		case "-armor":
			if i+1 < len(args) {
				input.Armor = args[i+1]
				i++
			}
		case "-shield":
			if i+1 < len(args) {
				input.Shield = args[i+1]
				i++
			}
		case "-slot":
			if i+1 < len(args) {
				input.Slot = args[i+1]
				i++
			}
		}
	}

	return input
}
