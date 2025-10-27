package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"ddsheetfinal/helpers"
)

const (
	slotMainHand = "main hand"
	slotOffHand  = "off hand"
)

func normalizeSlot(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "main", "main hand", "main-hand", "mainhand":
		return slotMainHand
	case "off", "off hand", "off-hand", "offhand":
		return slotOffHand
	default:
		return s
	}
}

func readCharacter(name string) (helpers.Character, error) {
	filename := helpers.CharacterPath(name)
	data, err := os.ReadFile(filename)
	if err != nil {
		return helpers.Character{}, fmt.Errorf("could not read character file: %v", err)
	}
	var char helpers.Character
	if err := json.Unmarshal(data, &char); err != nil {
		return helpers.Character{}, fmt.Errorf("could not parse character data: %v", err)
	}
	return char, nil
}

func equipWeapon(char *helpers.Character, weapon, rawSlot string) error {
	targetSlot := normalizeSlot(rawSlot)
	if targetSlot == "" {
		targetSlot = slotMainHand
	}
	if targetSlot != slotMainHand && targetSlot != slotOffHand {
		return fmt.Errorf("invalid slot: %s", rawSlot)
	}
	if targetSlot == slotMainHand && char.Weapon != "" {
		return fmt.Errorf("%s already occupied", slotMainHand)
	}
	if targetSlot == slotOffHand && char.OffHand != "" {
		return fmt.Errorf("%s already occupied", slotOffHand)
	}
	if targetSlot == slotMainHand {
		char.Weapon = weapon
	} else {
		char.OffHand = weapon
	}
	fmt.Printf("Equipped weapon %s to %s\n", weapon, targetSlot)
	return nil
}

func equipArmor(char *helpers.Character, armor string) error {
	if char.Armor != "" {
		return fmt.Errorf("armor slot already occupied")
	}
	char.Armor = armor
	fmt.Printf("Equipped armor %s\n", armor)
	return nil
}

func equipShield(char *helpers.Character, shield string) error {
	if char.Shield != "" {
		return fmt.Errorf("shield slot already occupied")
	}
	char.Shield = shield
	fmt.Printf("Equipped shield %s\n", shield)
	return nil
}

func parseEquipFlags(args []string) (name, weapon, armor, shield, slot string, parseErr bool) {
	equipCmd := flag.NewFlagSet("equip", flag.ExitOnError)
	equipCmd.Usage = func() {
		/* Intentionally left empty: override the default flag.Usage to suppress
		   automatic usage output from the flag package because this command prints
		   its own error messages; providing a no-op prevents unexpected stdout/stderr
		   output while keeping flag parsing behavior unchanged. */
	}
	namePtr := equipCmd.String("name", "", "character name (required)")
	weaponPtr := equipCmd.String("weapon", "", "weapon to equip")
	armorPtr := equipCmd.String("armor", "", "armor to equip")
	shieldPtr := equipCmd.String("shield", "", "shield to equip")
	slotPtr := equipCmd.String("slot", "", "slot to equip (e.g. \"main hand\", \"off hand\")")

	if err := equipCmd.Parse(args); err != nil {
		fmt.Println("error parsing flags")
		return "", "", "", "", "", true
	}
	if *namePtr == "" {
		fmt.Println("name is required")
		return "", "", "", "", "", true
	}
	return *namePtr, *weaponPtr, *armorPtr, *shieldPtr, *slotPtr, false
}

func validateItems(weapon, armor, shield string) bool {
	if weapon != "" && !helpers.IsValidEquipment(weapon) {
		fmt.Printf("Weapon '%s' not found in 5e-SRD-Equipment.csv\n", weapon)
		return false
	}
	if armor != "" && !helpers.IsValidEquipment(armor) {
		fmt.Printf("Armor '%s' not found in 5e-SRD-Equipment.csv\n", armor)
		return false
	}
	if shield != "" && !helpers.IsValidEquipment(shield) {
		fmt.Printf("Shield '%s' not found in 5e-SRD-Equipment.csv\n", shield)
		return false
	}
	return true
}

func applyEquipment(char *helpers.Character, weapon, armor, shield, slot string) (bool, error) {
	changed := false

	if weapon != "" {
		if err := equipWeapon(char, weapon, slot); err != nil {
			return false, err
		}
		changed = true
	}

	if armor != "" {
		if err := equipArmor(char, armor); err != nil {
			return false, err
		}
		changed = true
	}

	if shield != "" {
		if err := equipShield(char, shield); err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}

func ExecuteEquip(args []string) {
	name, weapon, armor, shield, slot, parseErr := parseEquipFlags(args)
	if parseErr {
		os.Exit(2)
	}

	if !validateItems(weapon, armor, shield) {
		os.Exit(2)
	}

	char, err := readCharacter(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	changed, err := applyEquipment(&char, weapon, armor, shield, slot)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !changed {
		fmt.Println("No equipment specified to equip.")
		os.Exit(1)
	}

	if err := helpers.SaveCharacter(char); err != nil {
		fmt.Println("Error saving character:", err)
		os.Exit(1)
	}
}
