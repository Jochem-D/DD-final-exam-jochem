package cli

import (
	"fmt"
)

// HelpHandler handles the help command
type HelpHandler struct{}

// NewHelpHandler creates a new help handler
func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}

// Handle processes the help command
func (h *HelpHandler) Handle(args []string) {
	fmt.Print(`Available commands:

  create
    Create a new character
    Usage: create -name NAME -race RACE -class CLASS -level N -str N -dex N -con N -int N -wis N -cha N

  view
    View a character's details
    Usage: view -name NAME

  list
    List all characters
    Usage: list

  delete
    Delete a character
    Usage: delete -name NAME

  equip
    Equip weapon, armor, or shield
    Usage:
      equip -name NAME -weapon WEAPON_NAME
      equip -name NAME -armor ARMOR_NAME
      equip -name NAME -shield SHIELD_NAME

  unequip
    Unequip weapon, armor, or shield (moves to inventory)
    Usage:
      unequip -name NAME -weapon
      unequip -name NAME -armor
      unequip -name NAME -shield

  learn-spell
    Learn a new spell
    Usage: learn-spell -name NAME -spell SPELL_NAME

  prepare-spell
    Prepare a learned spell
    Usage: prepare-spell -name NAME -spell SPELL_NAME

  learnable-spells
    List spells the character can still learn
    Usage: learnable-spells -name NAME

  help
    Show this help message

  serve
    Run a local HTTP server to serve the frontend and API endpoints
    Usage: serve [-dir PATH] [-port PORT]
`)
}
