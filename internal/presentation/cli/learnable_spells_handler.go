package cli

import (
	"fmt"
	"os"
	"strings"

	"ddsheetfinal/internal/application/usecases"
)

// LearnableSpellsHandler handles the learnable-spells command
type LearnableSpellsHandler struct {
	getLearnableSpellsUseCase *usecases.GetLearnableSpellsUseCase
}

// NewLearnableSpellsHandler creates a new learnable spells handler
func NewLearnableSpellsHandler(getLearnableSpellsUseCase *usecases.GetLearnableSpellsUseCase) *LearnableSpellsHandler {
	return &LearnableSpellsHandler{
		getLearnableSpellsUseCase: getLearnableSpellsUseCase,
	}
}

// Handle processes the learnable-spells command
func (h *LearnableSpellsHandler) Handle(args []string) {
	// Parse name
	name := h.parseName(args)
	if name == "" {
		fmt.Println("name is required")
		os.Exit(1)
	}

	// Execute use case
	spells, err := h.getLearnableSpellsUseCase.Execute(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Display results
	if len(spells) == 0 {
		fmt.Println("No learnable spells found")
	} else {
		fmt.Printf("Learnable spells: %s\n", strings.Join(spells, ", "))
	}
}

func (h *LearnableSpellsHandler) parseName(args []string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == "-name" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
