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
	input := h.getDefaultCharacter()

	for i := 0; i < len(args); i++ {
		if i+1 < len(args) {
			i = h.parseFlag(args[i], args[i+1], &input, i)
		}
	}

	return h.validateInput(input)
}

func (h *CreateHandler) getDefaultCharacter() dtos.CreateCharacterDTO {
	return dtos.CreateCharacterDTO{
		Level:      1,
		Str:        10,
		Dex:        10,
		Con:        10,
		Int:        10,
		Wis:        10,
		Cha:        10,
		Background: "acolyte",
	}
}

func (h *CreateHandler) parseFlag(flag, value string, input *dtos.CreateCharacterDTO, currentIndex int) int {
	switch flag {
	case "-name":
		input.Name = value
		return currentIndex + 1
	case "-race":
		input.Race = value
		return currentIndex + 1
	case "-class":
		input.Class = value
		return currentIndex + 1
	case "-background":
		input.Background = value
		return currentIndex + 1
	case "-level":
		fmt.Sscanf(value, "%d", &input.Level)
		return currentIndex + 1
	case "-str":
		fmt.Sscanf(value, "%d", &input.Str)
		return currentIndex + 1
	case "-dex":
		fmt.Sscanf(value, "%d", &input.Dex)
		return currentIndex + 1
	case "-con":
		fmt.Sscanf(value, "%d", &input.Con)
		return currentIndex + 1
	case "-int":
		fmt.Sscanf(value, "%d", &input.Int)
		return currentIndex + 1
	case "-wis":
		fmt.Sscanf(value, "%d", &input.Wis)
		return currentIndex + 1
	case "-cha":
		fmt.Sscanf(value, "%d", &input.Cha)
		return currentIndex + 1
	case "-skills":
		h.parseSkills(value, input)
		return currentIndex + 1
	}
	return currentIndex
}

func (h *CreateHandler) parseSkills(skillsStr string, input *dtos.CreateCharacterDTO) {
	input.Skills = strings.Split(skillsStr, ",")
	for j := range input.Skills {
		input.Skills[j] = strings.TrimSpace(input.Skills[j])
	}
}

func (h *CreateHandler) validateInput(input dtos.CreateCharacterDTO) (dtos.CreateCharacterDTO, error) {
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
