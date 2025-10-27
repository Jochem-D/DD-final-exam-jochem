package valueobjects

// ClassSkillProficiencies maps class names to their available skill proficiencies
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

// BackgroundSkillProficiencies maps background names to their skill proficiencies
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

// ClassSkillChoices maps class names to the number of skill choices they get
var ClassSkillChoices = map[string]int{
	"barbarian": 2, "bard": 3, "cleric": 2, "druid": 2, "fighter": 2, "monk": 2,
	"paladin": 2, "ranger": 3, "rogue": 4, "sorcerer": 2, "warlock": 2, "wizard": 2,
}

// RacialAbilityBonuses maps race names to their ability score bonuses
var RacialAbilityBonuses = map[string]map[string]int{
	"dwarf":              {"con": 2},
	"hill dwarf":         {"con": 2, "wis": 1},
	"mountain dwarf":     {"con": 2, "str": 2},
	"elf":                {"dex": 2},
	"high elf":           {"dex": 2, "int": 1},
	"wood elf":           {"dex": 2, "wis": 1},
	"dark elf":           {"dex": 2, "cha": 1},
	"halfling":           {"dex": 2},
	"lightfoot halfling": {"dex": 2, "cha": 1},
	"stout halfling":     {"dex": 2, "con": 1},
	"human":              {"str": 1, "dex": 1, "con": 1, "int": 1, "wis": 1, "cha": 1},
	"dragonborn":         {"str": 2, "cha": 1},
	"gnome":              {"int": 2},
	"forest gnome":       {"int": 2, "dex": 1},
	"rock gnome":         {"int": 2, "con": 1},
	"half-elf":           {"cha": 2, "choice1": 1, "choice2": 1},
	"half elf":           {"cha": 2, "choice1": 1, "choice2": 1},
	"half-orc":           {"str": 2, "con": 1},
	"half orc":           {"str": 2, "con": 1},
	"tiefling":           {"cha": 2, "int": 1},
}

// ClassSavingThrowProficiencies maps class names to their saving throw proficiencies
var ClassSavingThrowProficiencies = map[string][]string{
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

// SkillToAbility maps skill names to their governing ability
var SkillToAbility = map[string]string{
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

// ClassCantrips maps class names to their starting cantrips
var ClassCantrips = map[string][]string{
	"cleric":   {"guidance", "light"},
	"wizard":   {"mage hand", "fire bolt"},
	"druid":    {"shillelagh", "produce flame"},
	"bard":     {"vicious mockery", "mage hand"},
	"sorcerer": {"fire bolt", "prestidigitation"},
	"warlock":  {"eldritch blast"},
	"paladin":  {},
	"ranger":   {},
}
