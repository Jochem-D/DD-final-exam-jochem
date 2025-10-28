package valueobjects

import "strings"

// ArmorType represents the type of armor (light, medium, heavy)
type ArmorType string

const (
	ArmorTypeLight  ArmorType = "light"
	ArmorTypeMedium ArmorType = "medium"
	ArmorTypeHeavy  ArmorType = "heavy"
)

// ArmorInfo contains information about a specific armor piece
type ArmorInfo struct {
	Name string
	Base int
	Type ArmorType
}

// GetArmorInfo returns armor information for a given armor name (case-insensitive)
// Handles both "leather" and "leather armor" formats automatically
func GetArmorInfo(name string) (ArmorInfo, bool) {
	// Normalize the name: lowercase and try with/without " armor" suffix
	name = strings.ToLower(strings.TrimSpace(name))
	
	// Base armor data (without duplicates)
	armorData := map[string]ArmorInfo{
		// Light armor
		"padded":          {Name: "Padded", Base: 11, Type: ArmorTypeLight},
		"leather":         {Name: "Leather", Base: 11, Type: ArmorTypeLight},
		"studded leather": {Name: "Studded Leather", Base: 12, Type: ArmorTypeLight},

		// Medium armor
		"hide":        {Name: "Hide", Base: 12, Type: ArmorTypeMedium},
		"chain shirt": {Name: "Chain Shirt", Base: 13, Type: ArmorTypeMedium},
		"scale mail":  {Name: "Scale Mail", Base: 14, Type: ArmorTypeMedium},
		"breastplate": {Name: "Breastplate", Base: 14, Type: ArmorTypeMedium},
		"half plate":  {Name: "Half Plate", Base: 15, Type: ArmorTypeMedium},

		// Heavy armor
		"ring mail":  {Name: "Ring Mail", Base: 14, Type: ArmorTypeHeavy},
		"chain mail": {Name: "Chain Mail", Base: 16, Type: ArmorTypeHeavy},
		"splint":     {Name: "Splint", Base: 17, Type: ArmorTypeHeavy},
		"plate":      {Name: "Plate", Base: 18, Type: ArmorTypeHeavy},
	}

	// Try exact match first
	if info, ok := armorData[name]; ok {
		return info, true
	}

	// Try removing " armor" suffix
	if strings.HasSuffix(name, " armor") {
		baseName := strings.TrimSuffix(name, " armor")
		if info, ok := armorData[baseName]; ok {
			return info, true
		}
	}

	return ArmorInfo{}, false
}

