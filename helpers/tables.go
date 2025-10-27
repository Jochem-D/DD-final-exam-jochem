// helpers/tables.go
package helpers

// armorEntry models AC base and category for armor rules used by ComputeAC.
type armorEntry struct {
	Base int    // base AC
	Type string // "light" | "medium" | "heavy"
}

// armorTable maps lowercased armor names to their base and type.
// Source: 5e (PHB 2014) defaults.
var armorTable = map[string]armorEntry{
	// Light
	"padded":          {Base: 11, Type: "light"},
	"leather":         {Base: 11, Type: "light"},
	"studded leather": {Base: 12, Type: "light"},

	// Medium
	"hide":        {Base: 12, Type: "medium"},
	"chain shirt": {Base: 13, Type: "medium"},
	"scale mail":  {Base: 14, Type: "medium"},
	"breastplate": {Base: 14, Type: "medium"},
	"half plate":  {Base: 15, Type: "medium"},

	// Heavy
	"ring mail":  {Base: 14, Type: "heavy"},
	"chain mail": {Base: 16, Type: "heavy"},
	"splint":     {Base: 17, Type: "heavy"},
	"plate":      {Base: 18, Type: "heavy"},
}
