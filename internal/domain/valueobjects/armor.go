package valueobjects

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
func GetArmorInfo(name string) (ArmorInfo, bool) {
	armorData := map[string]ArmorInfo{
		// Light
		"padded":          {Name: "Padded", Base: 11, Type: ArmorTypeLight},
		"leather":         {Name: "Leather", Base: 11, Type: ArmorTypeLight},
		"studded leather": {Name: "Studded Leather", Base: 12, Type: ArmorTypeLight},

		// Medium
		"hide":        {Name: "Hide", Base: 12, Type: ArmorTypeMedium},
		"chain shirt": {Name: "Chain Shirt", Base: 13, Type: ArmorTypeMedium},
		"scale mail":  {Name: "Scale Mail", Base: 14, Type: ArmorTypeMedium},
		"breastplate": {Name: "Breastplate", Base: 14, Type: ArmorTypeMedium},
		"half plate":  {Name: "Half Plate", Base: 15, Type: ArmorTypeMedium},

		// Heavy
		"ring mail":  {Name: "Ring Mail", Base: 14, Type: ArmorTypeHeavy},
		"chain mail": {Name: "Chain Mail", Base: 16, Type: ArmorTypeHeavy},
		"splint":     {Name: "Splint", Base: 17, Type: ArmorTypeHeavy},
		"plate":      {Name: "Plate", Base: 18, Type: ArmorTypeHeavy},
	}

	info, ok := armorData[name]
	return info, ok
}
