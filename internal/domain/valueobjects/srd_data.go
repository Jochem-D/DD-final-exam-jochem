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

// SpellSlots returns spell slots per level for full casters (indexed by character level 1-20, then spell level 1-9)
// Based on D&D 5e spell slot progression table
var FullCasterSpellSlots = [][]int{
	//Lvl: 1  2  3  4  5  6  7  8  9  (spell levels)
	{2, 0, 0, 0, 0, 0, 0, 0, 0}, // Character level 1
	{3, 0, 0, 0, 0, 0, 0, 0, 0}, // 2
	{4, 2, 0, 0, 0, 0, 0, 0, 0}, // 3
	{4, 3, 0, 0, 0, 0, 0, 0, 0}, // 4
	{4, 3, 2, 0, 0, 0, 0, 0, 0}, // 5
	{4, 3, 3, 0, 0, 0, 0, 0, 0}, // 6
	{4, 3, 3, 1, 0, 0, 0, 0, 0}, // 7
	{4, 3, 3, 2, 0, 0, 0, 0, 0}, // 8
	{4, 3, 3, 3, 1, 0, 0, 0, 0}, // 9
	{4, 3, 3, 3, 2, 0, 0, 0, 0}, // 10
	{4, 3, 3, 3, 2, 1, 0, 0, 0}, // 11
	{4, 3, 3, 3, 2, 1, 0, 0, 0}, // 12
	{4, 3, 3, 3, 2, 1, 1, 0, 0}, // 13
	{4, 3, 3, 3, 2, 1, 1, 0, 0}, // 14
	{4, 3, 3, 3, 2, 1, 1, 1, 0}, // 15
	{4, 3, 3, 3, 2, 1, 1, 1, 0}, // 16
	{4, 3, 3, 3, 2, 1, 1, 1, 1}, // 17
	{4, 3, 3, 3, 3, 1, 1, 1, 1}, // 18
	{4, 3, 3, 3, 3, 2, 1, 1, 1}, // 19
	{4, 3, 3, 3, 3, 2, 2, 1, 1}, // 20
}

// HalfCasterSpellSlots for Paladin and Ranger
var HalfCasterSpellSlots = [][]int{
	//Lvl: 1  2  3  4  5  (spell levels - half casters only go to 5)
	{0, 0, 0, 0, 0}, // Character level 1
	{2, 0, 0, 0, 0}, // 2
	{3, 0, 0, 0, 0}, // 3
	{3, 0, 0, 0, 0}, // 4
	{4, 2, 0, 0, 0}, // 5
	{4, 2, 0, 0, 0}, // 6
	{4, 3, 0, 0, 0}, // 7
	{4, 3, 0, 0, 0}, // 8
	{4, 3, 2, 0, 0}, // 9
	{4, 3, 2, 0, 0}, // 10
	{4, 3, 3, 0, 0}, // 11
	{4, 3, 3, 0, 0}, // 12
	{4, 3, 3, 1, 0}, // 13
	{4, 3, 3, 1, 0}, // 14
	{4, 3, 3, 2, 0}, // 15
	{4, 3, 3, 2, 0}, // 16
	{4, 3, 3, 3, 1}, // 17
	{4, 3, 3, 3, 1}, // 18
	{4, 3, 3, 3, 2}, // 19
	{4, 3, 3, 3, 2}, // 20
}

// GetSpellSlots returns the spell slots for a class at a given level
// Returns nil if the class doesn't cast spells
// Index 0 is for cantrips, indices 1-9 are for spell levels 1-9
func GetSpellSlots(class string, level int) []int {
	if level < 1 || level > 20 {
		return nil
	}

	classLower := class
	if class != "" {
		classLower = class
	}

	// Warlock has special pact magic - check first before full caster logic
	if classLower == "warlock" {
		return getWarlockSpellSlots(level)
	}

	// Use centralized spellcasting logic
	isFullCaster := IsSpellcaster(classLower) && !IsPreparedCaster(classLower) || 
		classLower == "bard" || classLower == "sorcerer" || classLower == "wizard" || 
		classLower == "cleric" || classLower == "druid"
	
	isHalfCaster := classLower == "paladin" || classLower == "ranger"

	if isFullCaster {
		slots := make([]int, 10) // 0 for cantrips, 1-9 for spell levels
		cantrips := getCantripsKnown(classLower, level)
		if cantrips > 0 {
			slots[0] = cantrips
		}
		// Copy spell slots (skip index 0)
		spellSlots := FullCasterSpellSlots[level-1]
		copy(slots[1:], spellSlots)
		return slots
	}
	
	if isHalfCaster {
		slots := make([]int, 10)
		// Half casters don't get cantrips
		spellSlots := HalfCasterSpellSlots[level-1]
		copy(slots[1:], spellSlots)
		return slots
	}

	return nil
}

// getCantripsKnown returns the number of cantrips known for full casters at a given level
func getCantripsKnown(class string, level int) int {
	// Simplified - most full casters get 4 cantrips at level 1, more at higher levels
	switch class {
	case "bard", "sorcerer", "warlock":
		if level >= 10 {
			return 4
		}
		if level >= 4 {
			return 3
		}
		return 2
	case "cleric", "druid", "wizard":
		if level >= 10 {
			return 5
		}
		if level >= 4 {
			return 4
		}
		return 3
	}
	return 0
}

// getWarlockSpellSlots returns warlock pact magic slots
// Index 0 is for cantrips, then spell level slots
func getWarlockSpellSlots(level int) []int {
	slots := make([]int, 10) // 0 for cantrips, 1-9 for spell levels
	
	// Warlocks get cantrips
	cantrips := 2
	if level >= 4 {
		cantrips = 3
	}
	if level >= 10 {
		cantrips = 4
	}
	slots[0] = cantrips
	
	// Warlocks get 1 slot at level 1, 2 at level 2, 3 at level 11, 4 at level 17
	numSlots := 1
	if level >= 2 {
		numSlots = 2
	}
	if level >= 11 {
		numSlots = 3
	}
	if level >= 17 {
		numSlots = 4
	}
	
	// Slot level increases with warlock level
	slotLevel := 1
	if level >= 3 {
		slotLevel = 2
	}
	if level >= 5 {
		slotLevel = 3
	}
	if level >= 7 {
		slotLevel = 4
	}
	if level >= 9 {
		slotLevel = 5
	}
	
	// Put all slots at the appropriate level (1-indexed in the array at position slotLevel)
	slots[slotLevel] = numSlots
	
	return slots
}
