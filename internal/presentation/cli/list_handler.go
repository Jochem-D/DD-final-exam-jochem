package cli

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/application/usecases"
)

// ListHandler handles the list command
type ListHandler struct {
	listCharactersUseCase *usecases.ListCharactersUseCase
}

// NewListHandler creates a new list handler
func NewListHandler(listCharactersUseCase *usecases.ListCharactersUseCase) *ListHandler {
	return &ListHandler{
		listCharactersUseCase: listCharactersUseCase,
	}
}

// Handle processes the list command
func (h *ListHandler) Handle(args []string) {
	// Execute use case
	result, err := h.listCharactersUseCase.Execute()
	if err != nil {
		fmt.Println("Error listing characters:", err)
		os.Exit(1)
	}

	// Display results
	fmt.Println("Characters:")
	if len(result.Names) == 0 {
		fmt.Println("  (none found)")
	} else {
		for _, name := range result.Names {
			fmt.Printf(" - %s\n", name)
		}
	}
}
