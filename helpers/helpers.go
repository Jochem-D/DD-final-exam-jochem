package helpers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Character struct
type Character struct {
	Name       string `json:"name"`
	Race       string `json:"race"`
	Class      string `json:"class"`
	Background string `json:"background,omitempty"`
	Level      int    `json:"level"`
	Str        int    `json:"str"`
	Dex        int    `json:"dex"`
	Con        int    `json:"con"`
	Int        int    `json:"int"`
	Wis        int    `json:"wis"`
	Cha        int    `json:"cha"`

	ProficiencyBonus   int      `json:"proficiency_bonus,omitempty"`
	SkillProficiencies []string `json:"skill_proficiencies,omitempty"`

	Weapon  string `json:"weapon,omitempty"`
	OffHand string `json:"off_hand,omitempty"`
	Armor   string `json:"armor,omitempty"`
	Shield  string `json:"shield,omitempty"`

	Inventory []string `json:"inventory,omitempty"`

	Spells         []string `json:"spells,omitempty"`
	PreparedSpells []string `json:"prepared_spells,omitempty"`
}

// characterDir returns the directory for character files, creating it if it doesn't exist.
func characterDir() string {
	dir := "characters"
	// make sure the folder exists
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// CharacterPath returns the file path for a character's JSON file.
func CharacterPath(name string) string {
	return filepath.Join(characterDir(), fmt.Sprintf("%s.json", name))
}

// Save character to JSON
func SaveCharacter(c Character) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	filename := CharacterPath(c.Name)
	return os.WriteFile(filename, data, 0o644)
}



// ---------------- Combat / derived stats ----------------

func HasSkill(c *Character, skill string) bool {
	skill = strings.ToLower(skill)
	for _, s := range c.SkillProficiencies {
		if strings.ToLower(s) == skill {
			return true
		}
	}
	return false
}

func ComputeArmorClass(c *Character) int {
	dexMod := AbilityMod(c.Dex)
	conMod := AbilityMod(c.Con)
	wisMod := AbilityMod(c.Wis)

	ac := 10 + dexMod // default
	armorName := strings.ToLower(strings.TrimSpace(c.Armor))

	if armorName != "" {
		if base, kind, ok := lookupArmorEntry(armorName); ok {
			ac = base
			switch kind {
			case "light":
				ac += dexMod
			case "medium":
				if dexMod > 2 {
					ac += 2
				} else if dexMod > 0 {
					ac += dexMod
				}
			case "heavy":
				// no Dex
			}
		}
	} else {
		switch strings.ToLower(c.Class) {
		case "barbarian":
			ac = 10 + dexMod + conMod
		case "monk":
			ac = 10 + dexMod + wisMod
		}
	}

	if strings.TrimSpace(c.Shield) != "" {
		ac += 2
	}
	return ac
}

func InitiativeBonus(c *Character) int {
	return AbilityMod(c.Dex)
}

func PassivePerception(c *Character) int {
	// keep pointer signature for existing callers, but delegate to SRD-correct calc
	return PassivePerceptionRAW(*c)
}

// ---------------- Existing data tables ----------------

var ClassSkillProficiencies = map[string][]string{
	"barbarian": {"animal handling", "athletics", "intimidation", "nature", "perception", "survival"},
	"bard":      {"acrobatics", "animal handling", "arcana", "athletics", "deception", "history", "insight", "intimidation", "investigation", "medicine", "nature", "perception", "performance", "persuasion", "religion", "sleight of hand", "stealth", "survival"},
	"cleric":    {"history", "insight", "medicine", "persuasion", "religion"},
	"druid":     {"arcana", "animal handling", "insight", "medicine", "nature", "perception", "religion", "survival"},
	"fighter":   {"acrobatics", "animal handling", "athletics", "history", "insight", "intimidation", "perception", "survival"},
	"monk":      {"acrobatics", "athletics", "history", "insight", "religion", "stealth"},
	"paladin":   {"athletics", "insight", "intimidation", "medicine", "persuasion", "religion"},
	"ranger":    {"animal handling", "athletics", "insight", "investigation", "nature", "perception", "stealth", "survival"},
	"rogue":     {"acrobatics", "athletics", "deception", "insight", "intimidation", "investigation", "perception", "performance", "persuasion", "sleight of hand", "stealth"},
	"sorcerer":  {"arcana", "deception", "insight", "intimidation", "persuasion", "religion"},
	"warlock":   {"arcana", "deception", "history", "intimidation", "investigation", "nature", "religion"},
	"wizard":    {"arcana", "history", "insight", "investigation", "medicine", "religion"},
}

var BackgroundSkillProficiencies = map[string][]string{
	"acolyte":       {"insight", "religion"},
	"charlatan":     {"deception", "sleight of hand"},
	"criminal":      {"deception", "stealth"},
	"entertainer":   {"acrobatics", "performance"},
	"folk hero":     {"animal handling", "survival"},
	"guild artisan": {"insight", "persuasion"},
	"hermit":        {"medicine", "religion"},
	"noble":         {"history", "persuasion"},
	"outlander":     {"athletics", "survival"},
	"sage":          {"arcana", "history"},
	"sailor":        {"athletics", "perception"},
	"soldier":       {"athletics", "intimidation"},
	"urchin":        {"sleight of hand", "stealth"},
}

var ClassSkillChoices = map[string]int{
	"barbarian": 2, "bard": 3, "cleric": 2, "druid": 2, "fighter": 2, "monk": 2,
	"paladin": 2, "ranger": 3, "rogue": 4, "sorcerer": 2, "warlock": 2, "wizard": 2,
}

var RacialAbilityBonuses = map[string]map[string]int{
	"dwarf":              {"con": 2},
	"hill dwarf":         {"con": 2, "wis": 1},
	"mountain dwarf":     {"con": 2, "str": 2},
	"elf":                {"dex": 2},
	"high elf":           {"dex": 2, "int": 1},
	"wood elf":           {"dex": 2, "wis": 1},
	"dark elf":           {"dex": 2, "cha": 1}, // Drow
	"halfling":           {"dex": 2},
	"lightfoot halfling": {"dex": 2, "cha": 1},
	"stout halfling":     {"dex": 2, "con": 1},
	"human":              {"str": 1, "dex": 1, "con": 1, "int": 1, "wis": 1, "cha": 1},
	"dragonborn":         {"str": 2, "cha": 1},
	"gnome":              {"int": 2},
	"forest gnome":       {"int": 2, "dex": 1},
	"rock gnome":         {"int": 2, "con": 1},

	// Half-elf / Half-orc — include BOTH spellings (hyphen and space)
	"half-elf": {"cha": 2, "choice1": 1, "choice2": 1},
	"half elf": {"cha": 2, "choice1": 1, "choice2": 1},
	"half-orc": {"str": 2, "con": 1},
	"half orc": {"str": 2, "con": 1},

	"tiefling": {"cha": 2, "int": 1},
}

var classCantrips = map[string][]string{
	"cleric":   {"guidance", "light"},
	"wizard":   {"mage hand", "fire bolt"},
	"druid":    {"shillelagh", "produce flame"},
	"bard":     {"vicious mockery", "mage hand"},
	"sorcerer": {"fire bolt", "prestidigitation"},
	"warlock":  {"eldritch blast"},
	"paladin":  {},
	"ranger":   {},
}

// half-caster progression (paladin, ranger) -> map levels to full-caster level equivalence
var halfCasterLevelMap = map[int]int{
	1: 0, 2: 1, 3: 2, 4: 3, 5: 3, 6: 4, 7: 5, 8: 6, 9: 6, 10: 7,
	11: 8, 12: 9, 13: 9, 14: 10, 15: 11, 16: 12, 17: 12, 18: 13, 19: 13, 20: 14,
}

// ---------------- Name normalization / armor resolver (NEW) ----------------

const armorSuffix = " armor"

// canonEquipKey normalizes a name to a consistent lookup key.
func canonEquipKey(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.NewReplacer("’", "'", "–", "-", "—", "-", "/", " ").Replace(s)
	s = strings.Join(strings.Fields(s), " ") // collapse repeated spaces
	return s
}
// equipAlternates returns a list of alternate lookup keys for an equipment name
// (e.g. toggles the " armor" suffix and includes common aliases).
func equipAlternates(key string) []string {
	out := []string{key}
	// toggle " armor" suffix
	if strings.HasSuffix(key, armorSuffix) {
		out = append(out, strings.TrimSuffix(key, armorSuffix))
	} else {
		out = append(out, key+armorSuffix)
		out = append(out, key+" armor")
	}
	// common armor aliases
	switch key {
	case "leather":
		out = append(out, "leather armor")
	case "studded leather":
		out = append(out, "studded leather armor")
	case "scale mail":
		out = append(out, "scale mail armor")
	case "chain mail":
		out = append(out, "chain mail armor")
	case "ring mail":
		out = append(out, "ring mail armor")
	case "splint":
		out = append(out, "splint armor")
	case "plate":
		out = append(out, "plate armor")
	case "half plate":
		out = append(out, "half plate armor")
	}
	// de-duplicate
	seen := map[string]struct{}{}
	uniq := make([]string, 0, len(out))
	for _, k := range out {
		k = strings.TrimSpace(k)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		uniq = append(uniq, k)
	}
	return uniq
}

// lookupArmorEntry tries the exact key and alternates against armorTable.
func lookupArmorEntry(name string) (base int, kind string, ok bool) {
	key := canonEquipKey(name)
	if e, ok := armorTable[key]; ok {
		return e.Base, e.Type, true
	}
	for _, alt := range equipAlternates(key) {
		if e, ok := armorTable[alt]; ok {
			return e.Base, e.Type, true
		}
	}
	return 0, "", false
}

// ---------------- Simple CSV validators used by commands ----------------

func IsValidSpell(spell string) bool {
	spell = canonEquipKey(spell)
	f, err := os.Open("5e-SRD-Spells.csv")
	if err != nil {
		return false
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return false
	}
	for _, rec := range records {
		if len(rec) == 0 {
			continue
		}
		name := canonEquipKey(rec[0])
		if name == spell {
			return true
		}
	}
	return false
}

func IsValidEquipment(item string) bool {
	key := canonEquipKey(item)
	cands := equipAlternates(key)

	f, err := os.Open("5e-SRD-Equipment.csv")
	if err != nil {
		return false
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return false
	}
	for _, rec := range records {
		if len(rec) == 0 {
			continue
		}
		name := canonEquipKey(rec[0])
		// direct match
		if name == key {
			return true
		}
		// match any candidate
		if slices.Contains(cands, name) {
			return true
		}
	}
	return false
}

// Silence unused variable errors (remove these lines once you actually use the maps).
var (
	_ = classCantrips
	_ = halfCasterLevelMap
)
