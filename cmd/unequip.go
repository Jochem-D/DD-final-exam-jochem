package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"ddsheetfinal/helpers"
)

func ExecuteUnequip(args []string) {
	unequipCmd := flag.NewFlagSet("unequip", flag.ExitOnError)
	name := unequipCmd.String("name", "", "character name (required)")
	weapon := unequipCmd.Bool("weapon", false, "unequip weapon")
	armor := unequipCmd.Bool("armor", false, "unequip armor")
	shield := unequipCmd.Bool("shield", false, "unequip shield")

	if err := unequipCmd.Parse(args); err != nil || *name == "" {
		fmt.Println("name is required")
		unequipCmd.Usage()
		os.Exit(2)
	}

	filename := fmt.Sprintf("%s.json", *name)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read character file: %v\n", err)
		os.Exit(1)
	}

	var char helpers.Character
	if err := json.Unmarshal(data, &char); err != nil {
		fmt.Printf("Could not parse character data: %v\n", err)
		os.Exit(1)
	}

	changed := false
	if *weapon && char.Weapon != "" {
		char.Inventory = append(char.Inventory, char.Weapon)
		fmt.Printf("Unequipped weapon: %s (added to inventory)\n", char.Weapon)
		char.Weapon = ""
		changed = true
	}
	if *armor && char.Armor != "" {
		char.Inventory = append(char.Inventory, char.Armor)
		fmt.Printf("Unequipped armor: %s (added to inventory)\n", char.Armor)
		char.Armor = ""
		changed = true
	}
	if *shield && char.Shield != "" {
		char.Inventory = append(char.Inventory, char.Shield)
		fmt.Printf("Unequipped shield: %s (added to inventory)\n", char.Shield)
		char.Shield = ""
		changed = true
	}

	if !changed {
		fmt.Println("No equipment was unequipped.")
		os.Exit(2)
	}

	if err := helpers.SaveCharacter(char); err != nil {
		fmt.Println("Error saving character:", err)
		os.Exit(1)
	}
	fmt.Println("Equipment updated successfully!")
}
