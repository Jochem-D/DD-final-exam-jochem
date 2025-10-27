package cli

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/application/usecases"
)

// DeleteHandler handles the delete command
type DeleteHandler struct {
	deleteCharacterUseCase *usecases.DeleteCharacterUseCase
}

// NewDeleteHandler creates a new delete handler
func NewDeleteHandler(deleteCharacterUseCase *usecases.DeleteCharacterUseCase) *DeleteHandler {
	return &DeleteHandler{
		deleteCharacterUseCase: deleteCharacterUseCase,
	}
}

// Handle processes the delete command
func (h *DeleteHandler) Handle(args []string) {
	// Parse name from args
	name := h.parseName(args)
	if name == "" {
		fmt.Println("name is required")
		os.Exit(2)
	}

	// Execute use case
	if err := h.deleteCharacterUseCase.Execute(name); err != nil {
		fmt.Printf("Could not delete character: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("deleted %s\n", name)
}

func (h *DeleteHandler) parseName(args []string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == "-name" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
