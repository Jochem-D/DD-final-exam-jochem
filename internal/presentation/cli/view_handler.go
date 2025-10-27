package cli

import (
	"fmt"
	"os"
	"strings"

	"ddsheetfinal/internal/application/usecases"
	"ddsheetfinal/internal/domain/valueobjects"
)

// ViewHandler handles the view command
type ViewHandler struct {
	viewCharacterUseCase *usecases.ViewCharacterUseCase
}

// NewViewHandler creates a new view handler
func NewViewHandler(viewCharacterUseCase *usecases.ViewCharacterUseCase) *ViewHandler {
	return &ViewHandler{
		viewCharacterUseCase: viewCharacterUseCase,
	}
}

// Handle processes the view command
func (h *ViewHandler) Handle(args []string) {
	// Parse name from args
	name := h.parseName(args)
	if name == "" {
		fmt.Println("name is required")
		os.Exit(2)
	}

	// Execute use case
	output, err := h.viewCharacterUseCase.Execute(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Display character
	h.displayCharacter(output)
}

func (h *ViewHandler) parseName(args []string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == "-name" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func (h *ViewHandler) displayCharacter(output *usecases.ViewCharacterOutput) {
	char := output.Character

	// Basic info
	fmt.Printf("Name: %s\n", char.Name)
	fmt.Printf("Class: %s\n", strings.ToLower(char.Class))
	fmt.Printf("Race: %s\n", strings.ToLower(char.Race))
	fmt.Printf("Background: %s\n", strings.ToLower(char.Background))
	fmt.Printf("Level: %d\n", char.Level)

	// Ability scores
	fmt.Println("Ability scores:")
	fmt.Printf("  STR: %d (%+d)\n", char.Str, valueobjects.AbilityModifier(char.Str))
	fmt.Printf("  DEX: %d (%+d)\n", char.Dex, valueobjects.AbilityModifier(char.Dex))
	fmt.Printf("  CON: %d (%+d)\n", char.Con, valueobjects.AbilityModifier(char.Con))
	fmt.Printf("  INT: %d (%+d)\n", char.Int, valueobjects.AbilityModifier(char.Int))
	fmt.Printf("  WIS: %d (%+d)\n", char.Wis, valueobjects.AbilityModifier(char.Wis))
	fmt.Printf("  CHA: %d (%+d)\n", char.Cha, valueobjects.AbilityModifier(char.Cha))
	fmt.Printf("Proficiency bonus: +%d\n", valueobjects.ProficiencyBonus(char.Level))

	// Skill proficiencies
	if len(char.SkillProficiencies) > 0 {
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(char.SkillProficiencies, ", "))
	}

	// Equipment
	if char.Weapon != "" || char.OffHand != "" || char.Armor != "" || char.Shield != "" {
		fmt.Println("Equipment:")
		if char.Weapon != "" {
			fmt.Printf("  Weapon: %s\n", char.Weapon)
		}
		if char.OffHand != "" {
			fmt.Printf("  Off-hand: %s\n", char.OffHand)
		}
		if char.Armor != "" {
			fmt.Printf("  Armor: %s\n", char.Armor)
		}
		if char.Shield != "" {
			fmt.Printf("  Shield: %s\n", char.Shield)
		}
	}

	// Derived stats
	fmt.Printf("Armor Class: %d (%s)\n", output.ArmorClass, output.ArmorClassDesc)
	fmt.Printf("Initiative: %+d\n", output.InitiativeBonus)
	fmt.Printf("Passive Perception: %d\n", output.PassivePerception)

	// Spells
	if len(char.Spells) > 0 {
		fmt.Printf("Known spells: %s\n", strings.Join(char.Spells, ", "))
	}
	if len(char.PreparedSpells) > 0 {
		fmt.Printf("Prepared spells: %s\n", strings.Join(char.PreparedSpells, ", "))
	}
}
