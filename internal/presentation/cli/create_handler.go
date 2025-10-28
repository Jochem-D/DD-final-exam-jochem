package cli

import (
	"fmt"
	"os"
	"strings"

	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/application/usecases"
)

// CreateHandler handles the create command
type CreateHandler struct {
	createCharacterUseCase *usecases.CreateCharacterUseCase
}

// NewCreateHandler creates a new create handler
func NewCreateHandler(createCharacterUseCase *usecases.CreateCharacterUseCase) *CreateHandler {
	return &CreateHandler{
		createCharacterUseCase: createCharacterUseCase,
	}
}

// Handle processes the create command
func (h *CreateHandler) Handle(args []string) {
	// Parse flags
	input, err := h.parseFlags(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Execute use case
	result, err := h.createCharacterUseCase.Execute(input)
	if err != nil {
		fmt.Println("Error creating character:", err)
		os.Exit(1)
	}

	fmt.Printf("saved character %s\n", result.Name)
}

func (h *CreateHandler) parseFlags(args []string) (dtos.CreateCharacterDTO, error) {
	// Simple flag parsing (in production, use flag package or a library)
	input := dtos.CreateCharacterDTO{
		Level:      1,
		Str:        10,
		Dex:        10,
		Con:        10,
		Int:        10,
		Wis:        10,
		Cha:        10,
		Background: "acolyte",
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-name":
			if i+1 < len(args) {
				input.Name = args[i+1]
				i++
			}
		case "-race":
			if i+1 < len(args) {
				input.Race = args[i+1]
				i++
			}
		case "-class":
			if i+1 < len(args) {
				input.Class = args[i+1]
				i++
			}
		case "-background":
			if i+1 < len(args) {
				input.Background = args[i+1]
				i++
			}
		case "-level":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Level)
				i++
			}
		case "-str":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Str)
				i++
			}
		case "-dex":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Dex)
				i++
			}
		case "-con":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Con)
				i++
			}
		case "-int":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Int)
				i++
			}
		case "-wis":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Wis)
				i++
			}
		case "-cha":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &input.Cha)
				i++
			}
		case "-skills":
			if i+1 < len(args) {
				skillsStr := args[i+1]
				input.Skills = strings.Split(skillsStr, ",")
				for j := range input.Skills {
					input.Skills[j] = strings.TrimSpace(input.Skills[j])
				}
				i++
			}
		}
	}

	// Validate required fields
	if input.Name == "" {
		return input, fmt.Errorf("name is required")
	}
	if input.Race == "" {
		return input, fmt.Errorf("race is required")
	}
	if input.Class == "" {
		return input, fmt.Errorf("class is required")
	}

	return input, nil
}
