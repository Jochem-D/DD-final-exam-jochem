package cli

import (
	"fmt"
	"os"
	"strings"

	"ddsheetfinal/internal/application/usecases"
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

	// Ability scores with modifiers from output
	fmt.Println("Ability scores:")
	fmt.Printf("  STR: %d (%+d)\n", char.Str, output.StrMod)
	fmt.Printf("  DEX: %d (%+d)\n", char.Dex, output.DexMod)
	fmt.Printf("  CON: %d (%+d)\n", char.Con, output.ConMod)
	fmt.Printf("  INT: %d (%+d)\n", char.Int, output.IntMod)
	fmt.Printf("  WIS: %d (%+d)\n", char.Wis, output.WisMod)
	fmt.Printf("  CHA: %d (%+d)\n", char.Cha, output.ChaMod)
	fmt.Printf("Proficiency bonus: +%d\n", output.ProficiencyBonus)

	// Skill proficiencies
	if len(char.SkillProficiencies) > 0 {
		fmt.Printf("Skill proficiencies: %s\n", strings.Join(char.SkillProficiencies, ", "))
	}

	// Spell slots (from output)
	if output.SpellSlots != nil {
		hasSlots := false
		for _, slots := range output.SpellSlots {
			if slots > 0 {
				hasSlots = true
				break
			}
		}
		if hasSlots {
			fmt.Println("Spell slots:")
			for i, slots := range output.SpellSlots {
				if slots > 0 {
					// Index 0 is cantrips (Level 0), rest are spell levels
					fmt.Printf("  Level %d: %d\n", i, slots)
				}
			}
			
			// Spellcasting information from output
			if output.SpellcastingAbility != "" {
				fmt.Printf("Spellcasting ability: %s\n", output.SpellcastingAbility)
				fmt.Printf("Spell save DC: %d\n", output.SpellSaveDC)
				fmt.Printf("Spell attack bonus: +%d\n", output.SpellAttackBonus)
			}
		}
	}

	// Equipment (no "Equipment:" header, just list items)
	if char.Weapon != "" {
		fmt.Printf("Main hand: %s\n", char.Weapon)
	}
	if char.OffHand != "" {
		fmt.Printf("Off hand: %s\n", char.OffHand)
	}
	if char.Armor != "" {
		fmt.Printf("Armor: %s\n", char.Armor)
	}
	if char.Shield != "" {
		fmt.Printf("Shield: %s\n", char.Shield)
	}

	// Derived stats
	fmt.Printf("Armor class: %d\n", output.ArmorClass)
	fmt.Printf("Initiative bonus: %d\n", output.InitiativeBonus)
	fmt.Printf("Passive perception: %d\n", output.PassivePerception)

	// Spells
	if len(char.Spells) > 0 {
		fmt.Printf("Known spells: %s\n", strings.Join(char.Spells, ", "))
	}
	if len(char.PreparedSpells) > 0 {
		fmt.Printf("Prepared spells: %s\n", strings.Join(char.PreparedSpells, ", "))
	}
}
