package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	fmtUnarmoredDex       = "10 + Dex(%+d)"
	fmtUnarmoredDexShield = "10 + Dex(%+d) + shield(+2)"
)

// ComputeAC calculates total Armor Class, including armor type rules and shields.
// Uses armorTable already defined in helpers.go.
// Rules:
// - Light:  base + Dex
// - Medium: base + min(Dex, +2)
// - Heavy:  base (ignore Dex)
// - Shield: +2 if equipped
// - Unarmored:
//   - Monk (no shield): 10 + Dex + Wis
//   - Barbarian:        10 + Dex + Con (+ shield if any)
//   - Otherwise:        10 + Dex (+ shield if any)
func ComputeAC(c Character) (int, string) {
	dexMod := AbilityMod(c.Dex)
	conMod := AbilityMod(c.Con)
	wisMod := AbilityMod(c.Wis)
	shieldBonus := 0
	if strings.TrimSpace(strings.ToLower(c.Shield)) != "" {
		shieldBonus = 2
	}
	armorName := strings.TrimSpace(c.Armor)

	if armorName == "" {
		return computeUnarmoredAC(c, dexMod, conMod, wisMod, shieldBonus)
	}

	if base, kind, ok := lookupArmorEntry(armorName); ok {
		return computeArmoredAC(kind, base, dexMod, shieldBonus)
	}
	return computeUnarmoredAC(c, dexMod, conMod, wisMod, shieldBonus)
}

func computeUnarmoredAC(c Character, dexMod, conMod, wisMod, shieldBonus int) (int, string) {
	classLower := strings.ToLower(strings.TrimSpace(c.Class))
	switch classLower {
	case "monk":
		if shieldBonus == 0 {
			ac := 10 + dexMod + wisMod
			return ac, fmt.Sprintf("10 + Dex(%+d) + Wis(%+d)", dexMod, wisMod)
		}
	case "barbarian":
		ac := 10 + dexMod + conMod + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("10 + Dex(%+d) + Con(%+d) + shield(+2)", dexMod, conMod)
		}
		return ac, fmt.Sprintf("10 + Dex(%+d) + Con(%+d)", dexMod, conMod)
	}
	ac := 10 + dexMod + shieldBonus
	if shieldBonus > 0 {
		return ac, fmt.Sprintf(fmtUnarmoredDexShield, dexMod)
	}
	return ac, fmt.Sprintf(fmtUnarmoredDex, dexMod)
}

func computeArmoredAC(kind string, base, dexMod, shieldBonus int) (int, string) {
	switch kind {
	case "light":
		ac := base + dexMod + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("%d (light base) + Dex(%+d) + shield(+2)", base, dexMod)
		}
		return ac, fmt.Sprintf("%d (light base) + Dex(%+d)", base, dexMod)
	case "medium":
		d := dexMod
		if d > 2 {
			d = 2
		}
		ac := base + d + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("%d (medium base) + Dex cap(+%d) + shield(+2)", base, d)
		}
		return ac, fmt.Sprintf("%d (medium base) + Dex cap(+%d)", base, d)
	case "heavy":
		ac := base + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf("%d (heavy base) + shield(+2)", base)
		}
		return ac, fmt.Sprintf("%d (heavy base)", base)
	default:
		ac := 10 + dexMod + shieldBonus
		if shieldBonus > 0 {
			return ac, fmt.Sprintf(fmtUnarmoredDexShield, dexMod)
		}
		return ac, fmt.Sprintf(fmtUnarmoredDex, dexMod)
	}
}
// --- derive API: compute derived numbers from a character JSON payload ---

type DerivedResponse struct {
	AbilityMods       map[string]int    `json:"ability_mods"`
	ProficiencyBonus  int               `json:"proficiency_bonus"`
	Saves             map[string]string `json:"saves,omitempty"`
	Skills            map[string]string `json:"skills,omitempty"`
	Initiative        int               `json:"initiative"`
	PassivePerception int               `json:"passive_perception"`
	ArmorClass        int               `json:"armor_class"`
	ACFormula         string            `json:"ac_formula,omitempty"`
}

// helper: format signed integer like "+3" or "-1"
func signedString(n int) string {
	if n >= 0 {
		return "+" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// DeriveFromMap computes derived fields from a generic JSON-decoded map.
// It tolerantly reads many common key variations to be robust with frontend data.
func DeriveFromMap(m map[string]any) DerivedResponse {
	str := readIntFromMap(m, "Str", "str", "Strength", "strength")
	dex := readIntFromMap(m, "Dex", "dex", "Dexterity", "dexterity")
	con := readIntFromMap(m, "Con", "con", "Constitution", "constitution")
	intel := readIntFromMap(m, "Int", "int", "Intelligence", "intelligence")
	wis := readIntFromMap(m, "Wis", "wis", "Wisdom", "wisdom")
	cha := readIntFromMap(m, "Cha", "cha", "Charisma", "charisma")
	level := readIntFromMap(m, "Level", "level")

	c := Character{
		Name:   readStringFromMap(m, "Name", "name"),
		Class:  readStringFromMap(m, "Class", "class"),
		Level:  level,
		Str:    str,
		Dex:    dex,
		Con:    con,
		Int:    intel,
		Wis:    wis,
		Cha:    cha,
		Armor:  readStringFromMap(m, "Armor", "armor"),
		Shield: readStringFromMap(m, "Shield", "shield"),
	}

	STR := AbilityMod(c.Str)
	DEX := AbilityMod(c.Dex)
	CON := AbilityMod(c.Con)
	INT := AbilityMod(c.Int)
	WIS := AbilityMod(c.Wis)
	CHA := AbilityMod(c.Cha)

	pb := GetProficiencyBonus(c.Level)

	mods := AbilityMods{
		Strength:     STR,
		Dexterity:    DEX,
		Constitution: CON,
		Intelligence: INT,
		Wisdom:       WIS,
		Charisma:     CHA,
	}
	saves := computeSaves(m, mods, pb)

	skillMap := map[string]string{
		"Acrobatics": "Dexterity", "Animal Handling": "Wisdom", "Arcana": "Intelligence",
		"Athletics": "Strength", "Deception": "Charisma", "History": "Intelligence",
		"Insight": "Wisdom", "Intimidation": "Charisma", "Investigation": "Intelligence",
		"Medicine": "Wisdom", "Nature": "Intelligence", "Perception": "Wisdom",
		"Performance": "Charisma", "Persuasion": "Charisma", "Religion": "Intelligence",
		"Sleight of Hand": "Dexterity", "Stealth": "Dexterity", "Survival": "Wisdom",
	}
	abilityByName := map[string]int{"Strength": STR, "Dexterity": DEX, "Constitution": CON, "Intelligence": INT, "Wisdom": WIS, "Charisma": CHA}
	profSkills := getProficiencySet(m, "SkillProficiencies", "skill_proficiencies")
	skills := computeSkills(skillMap, abilityByName, profSkills, pb)

	initVal := DEX
	pp := 10 + WIS

	hasShield := detectShield(m)
	if hasShield && strings.TrimSpace(c.Shield) == "" {
		c.Shield = "shield"
	}
	ac, formula := ComputeAC(c)

	return DerivedResponse{
		AbilityMods:       map[string]int{"Strength": STR, "Dexterity": DEX, "Constitution": CON, "Intelligence": INT, "Wisdom": WIS, "Charisma": CHA},
		ProficiencyBonus:  pb,
		Saves:             saves,
		Skills:            skills,
		Initiative:        initVal,
		PassivePerception: pp,
		ArmorClass:        ac,
		ACFormula:         formula,
	}
}

// Helper functions for DeriveFromMap

func readIntFromMap(m map[string]any, keys ...string) int {
	for _, k := range keys {
		if val, ok := m[k]; ok {
			if parsed, ok := parseInt(val); ok {
				return parsed
			}
		}
		lk := strings.ToLower(k)
		for mk, mv := range m {
			if strings.ToLower(mk) == lk {
				if parsed, ok := parseInt(mv); ok {
					return parsed
				}
			}
		}
	}
	return 0
}

// parseInt tries to parse various types to int, returns (value, true) if successful
func parseInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return int(i), true
		}
	case string:
		if i, err := strconv.Atoi(strings.TrimSpace(t)); err == nil {
			return i, true
		}
	}
	return 0, false
}

func readStringFromMap(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		lk := strings.ToLower(k)
		for mk, mv := range m {
			if strings.ToLower(mk) == lk {
				if s, ok := mv.(string); ok {
					return s
				}
			}
		}
	}
	return ""
}

func getProficiencySet(m map[string]any, keys ...string) map[string]bool {
	profSet := map[string]bool{}
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if arr, ok := v.([]any); ok {
				for _, e := range arr {
					if s, ok := e.(string); ok {
						profSet[strings.ToLower(strings.TrimSpace(s))] = true
					}
				}
			}
		}
	}
	return profSet
}

// AbilityMods is a lightweight struct used when computing saves.
type AbilityMods struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

func computeSaves(m map[string]any, mods AbilityMods, pb int) map[string]string {
	saveNames := []string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}
	abilityMap := map[string]int{
		"Strength":     mods.Strength,
		"Dexterity":    mods.Dexterity,
		"Constitution": mods.Constitution,
		"Intelligence": mods.Intelligence,
		"Wisdom":       mods.Wisdom,
		"Charisma":     mods.Charisma,
	}
	saves := make(map[string]string, len(saveNames))
	for _, s := range saveNames {
		prof := isSaveProficient(m, s)
		total := abilityMap[s]
		if prof {
			total += pb
		}
		saves[s] = signedString(total)
	}
	return saves
}

// isSaveProficient checks several key variants for a boolean proficiency flag.
func isSaveProficient(m map[string]any, save string) bool {
	keys := []string{save + "-save-prof", strings.ToLower(save) + "-save-prof"}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if b, ok := v.(bool); ok && b {
				return true
			}
		}
	}
	return false
}

func computeSkills(skillMap map[string]string, abilityByName map[string]int, profSkills map[string]bool, pb int) map[string]string {
	skills := map[string]string{}
	for k, ability := range skillMap {
		base := abilityByName[ability]
		total := base
		if profSkills[strings.ToLower(k)] {
			total += pb
		}
		skills[k] = signedString(total)
	}
	return skills
}

func detectShield(m map[string]any) bool {
	if v, ok := m["Equipment"]; ok {
		if arr, ok := v.([]any); ok {
			for _, e := range arr {
				if s, ok := e.(string); ok {
					if strings.Contains(strings.ToLower(s), "shield") {
						return true
					}
				}
			}
		}
	}
	if en, ok := m["enriched"]; ok {
		if em, ok := en.(map[string]any); ok {
			if eq, ok := em["equipment"]; ok {
				if eqm, ok := eq.(map[string]any); ok {
					for k := range eqm {
						if strings.Contains(strings.ToLower(k), "shield") {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// DeriveHandler HTTP endpoint: POST character JSON -> derived JSON
func DeriveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var m map[string]any
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}
	out := DeriveFromMap(m)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
