package cli

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/application/usecases"
)

// PrepareSpellHandler handles the prepare-spell command
type PrepareSpellHandler struct {
	prepareSpellUseCase *usecases.PrepareSpellUseCase
}

// NewPrepareSpellHandler creates a new prepare spell handler
func NewPrepareSpellHandler(prepareSpellUseCase *usecases.PrepareSpellUseCase) *PrepareSpellHandler {
	return &PrepareSpellHandler{
		prepareSpellUseCase: prepareSpellUseCase,
	}
}

// Handle processes the prepare-spell command
func (h *PrepareSpellHandler) Handle(args []string) {
	// Parse input
	input := h.parseInput(args)
	if input.CharacterName == "" || input.SpellName == "" {
		fmt.Println("name and spell are required")
		os.Exit(1)
	}

	// Execute use case
	if err := h.prepareSpellUseCase.Execute(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if input.Remove {
		fmt.Printf("Unprepared spell %s\n", input.SpellName)
	} else {
		fmt.Printf("Prepared spell %s\n", input.SpellName)
	}
}

func (h *PrepareSpellHandler) parseInput(args []string) usecases.PrepareSpellInput {
	input := usecases.PrepareSpellInput{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-name":
			if i+1 < len(args) {
				input.CharacterName = args[i+1]
				i++
			}
		case "-spell":
			if i+1 < len(args) {
				input.SpellName = args[i+1]
				i++
			}
		case "-remove":
			input.Remove = true
		}
	}

	return input
}
