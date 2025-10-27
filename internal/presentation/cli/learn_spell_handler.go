package cli

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/application/usecases"
)

// LearnSpellHandler handles the learn-spell command
type LearnSpellHandler struct {
	learnSpellUseCase *usecases.LearnSpellUseCase
}

// NewLearnSpellHandler creates a new learn spell handler
func NewLearnSpellHandler(learnSpellUseCase *usecases.LearnSpellUseCase) *LearnSpellHandler {
	return &LearnSpellHandler{
		learnSpellUseCase: learnSpellUseCase,
	}
}

// Handle processes the learn-spell command
func (h *LearnSpellHandler) Handle(args []string) {
	// Parse input
	name, spell := h.parseInput(args)
	if name == "" || spell == "" {
		fmt.Println("name and spell are required")
		os.Exit(1)
	}

	// Execute use case
	if err := h.learnSpellUseCase.Execute(name, spell); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Learned spell %s\n", spell)
}

func (h *LearnSpellHandler) parseInput(args []string) (name, spell string) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-name":
			if i+1 < len(args) {
				name = args[i+1]
				i++
			}
		case "-spell":
			if i+1 < len(args) {
				spell = args[i+1]
				i++
			}
		}
	}
	return
}
