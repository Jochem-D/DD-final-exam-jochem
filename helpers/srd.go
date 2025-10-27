package helpers

import "strings"

// --- internal helpers ---

func lower(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// normalize an ability name/abbr to 3-letter abbr: "strength" -> "str", "Str"->"str"
func abbr(ab string) string {
	switch lower(ab) {
	case "strength", "str":
		return "str"
	case "dexterity", "dex":
		return "dex"
	case "constitution", "con":
		return "con"
	case "intelligence", "int":
		return "int"
	case "wisdom", "wis":
		return "wis"
	case "charisma", "cha":
		return "cha"
	default:
		return lower(ab)
	}
}

// abilityScore gets a score by abbrev ("str","dex",...).
func abilityScore(c Character, ab string) int {
	switch abbr(ab) {
	case "str":
		return c.Str
	case "dex":
		return c.Dex
	case "con":
		return c.Con
	case "int":
		return c.Int
	case "wis":
		return c.Wis
	case "cha":
		return c.Cha
	default:
		return 10
	}
}

// --- SRD data (package-level, no per-call allocs) ---

// Class → saving throw proficiencies (ability abbreviations).
var classSaveProfs = map[string][]string{
	"barbarian": {"str", "con"},
	"bard":      {"dex", "cha"},
	"cleric":    {"wis", "cha"},
	"druid":     {"int", "wis"},
	"fighter":   {"str", "con"},
	"monk":      {"str", "dex"},
	"paladin":   {"wis", "cha"},
	"ranger":    {"str", "dex"},
	"rogue":     {"dex", "int"},
	"sorcerer":  {"con", "cha"},
	"warlock":   {"wis", "cha"},
	"wizard":    {"int", "wis"},
}

// Skill → governing ability (ability abbreviations).
var skillToAbility = map[string]string{
	"acrobatics":      "dex",
	"animal handling": "wis",
	"arcana":          "int",
	"athletics":       "str",
	"deception":       "cha",
	"history":         "int",
	"insight":         "wis",
	"intimidation":    "cha",
	"investigation":   "int",
	"medicine":        "wis",
	"nature":          "int",
	"perception":      "wis",
	"performance":     "cha",
	"persuasion":      "cha",
	"religion":        "int",
	"sleight of hand": "dex",
	"stealth":         "dex",
	"survival":        "wis",
}

// --- exported, DRY API ---

// IsSaveProficient reports if the character is proficient in a given save (by abbr or full name).
func IsSaveProficient(c Character, ab string) bool {
	ab = abbr(ab)
	profs := classSaveProfs[lower(c.Class)]
	for _, p := range profs {
		if p == ab {
			return true
		}
	}
	return false
}

// SavingThrowTotal = ability mod + proficiency bonus if proficient (SRD 5e).
func SavingThrowTotal(c Character, ab string) int {
	mod := AbilityMod(abilityScore(c, ab))
	if IsSaveProficient(c, ab) {
		return mod + GetProficiencyBonus(c.Level)
	}
	return mod
}

// SkillAbility returns the governing ability abbreviation ("str","dex",...) for a skill.
func SkillAbility(skill string) string {
	return skillToAbility[lower(skill)]
}

// SkillTotal = ability mod + proficiency bonus if the skill is in c.SkillProficiencies.
func SkillTotal(c Character, skill string) int {
	ability := SkillAbility(skill)
	mod := AbilityMod(abilityScore(c, ability))
	if HasSkill(&c, lower(skill)) { // reuse your existing helper (case-insensitive)
		return mod + GetProficiencyBonus(c.Level)
	}
	return mod
}

// PassivePerception (RAW) = 10 + Perception skill total.
func PassivePerceptionRAW(c Character) int {
	return 10 + SkillTotal(c, "perception")
}
