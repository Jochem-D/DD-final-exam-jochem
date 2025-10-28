package valueobjects

// Ability score constants
const (
	abilityStr = "str"
	abilityDex = "dex"
	abilityCon = "con"
	abilityInt = "int"
	abilityWis = "wis"
	abilityCha = "cha"
)

// Skill constants
const (
	skillAthletics      = "athletics"
	skillInsight        = "insight"
	skillReligion       = "religion"
	skillIntimidation   = "intimidation"
	skillHistory        = "history"
	skillSurvival       = "survival"
	skillPersuasion     = "persuasion"
	skillPerception     = "perception"
	skillStealth        = "stealth"
	skillMedicine       = "medicine"
	skillDeception      = "deception"
	skillArcana         = "arcana"
	skillAnimalHandling = "animal handling"
	skillSlightOfHand   = "sleight of hand"
)

// Class constants
const (
	classWizard = "wizard"
)

// ClassSkillProficiencies maps class names to their available skill proficiencies
var ClassSkillProficiencies = map[string][]string{
	"barbarian": {skillAnimalHandling, skillAthletics, skillIntimidation, "nature", skillPerception, skillSurvival},
	"bard":      {"acrobatics", skillAnimalHandling, skillArcana, skillAthletics, skillDeception, skillHistory, skillInsight, skillIntimidation, "investigation", skillMedicine, "nature", skillPerception, "performance", skillPersuasion, skillReligion, skillSlightOfHand, skillStealth, skillSurvival},
	"cleric":    {skillHistory, skillInsight, skillMedicine, skillPersuasion, skillReligion},
	"druid":     {skillArcana, skillAnimalHandling, skillInsight, skillMedicine, "nature", skillPerception, skillReligion, skillSurvival},
	"fighter":   {"acrobatics", skillAnimalHandling, skillAthletics, skillHistory, skillInsight, skillIntimidation, skillPerception, skillSurvival},
	"monk":      {"acrobatics", skillAthletics, skillHistory, skillInsight, skillReligion, skillStealth},
	"paladin":   {skillAthletics, skillInsight, skillIntimidation, skillMedicine, skillPersuasion, skillReligion},
	"ranger":    {skillAnimalHandling, skillAthletics, skillInsight, "investigation", "nature", skillPerception, skillStealth, skillSurvival},
	"rogue":     {"acrobatics", skillAthletics, skillDeception, skillInsight, skillIntimidation, "investigation", skillPerception, "performance", skillPersuasion, skillSlightOfHand, skillStealth},
	"sorcerer":  {skillArcana, skillDeception, skillInsight, skillIntimidation, skillPersuasion, skillReligion},
	"warlock":   {skillArcana, skillDeception, skillHistory, skillIntimidation, "investigation", "nature", skillReligion},
	classWizard: {skillArcana, skillHistory, skillInsight, "investigation", skillMedicine, skillReligion},
}

// BackgroundSkillProficiencies maps background names to their skill proficiencies
var BackgroundSkillProficiencies = map[string][]string{
	"acolyte":       {skillInsight, skillReligion},
	"charlatan":     {skillDeception, skillSlightOfHand},
	"criminal":      {skillDeception, skillStealth},
	"entertainer":   {"acrobatics", "performance"},
	"folk hero":     {skillAnimalHandling, skillSurvival},
	"guild artisan": {skillInsight, skillPersuasion},
	"hermit":        {skillMedicine, skillReligion},
	"noble":         {skillHistory, skillPersuasion},
	"outlander":     {skillAthletics, skillSurvival},
	"sage":          {skillArcana, skillHistory},
	"sailor":        {skillAthletics, skillPerception},
	"soldier":       {skillAthletics, skillIntimidation},
	"urchin":        {skillSlightOfHand, skillStealth},
}

// ClassSkillChoices maps class names to the number of skill choices they get
var ClassSkillChoices = map[string]int{
	"barbarian": 2, "bard": 3, "cleric": 2, "druid": 2, "fighter": 2, "monk": 2,
	"paladin": 2, "ranger": 3, "rogue": 4, "sorcerer": 2, "warlock": 2, "wizard": 2,
}

// RacialAbilityBonuses maps race names to their ability score bonuses
var RacialAbilityBonuses = map[string]map[string]int{
	"dwarf":              {abilityCon: 2},
	"hill dwarf":         {abilityCon: 2, abilityWis: 1},
	"mountain dwarf":     {abilityCon: 2, abilityStr: 2},
	"elf":                {abilityDex: 2},
	"high elf":           {abilityDex: 2, abilityInt: 1},
	"wood elf":           {abilityDex: 2, abilityWis: 1},
	"dark elf":           {abilityDex: 2, abilityCha: 1},
	"halfling":           {abilityDex: 2},
	"lightfoot halfling": {abilityDex: 2, abilityCha: 1},
	"stout halfling":     {abilityDex: 2, abilityCon: 1},
	"human":              {abilityStr: 1, abilityDex: 1, abilityCon: 1, abilityInt: 1, abilityWis: 1, abilityCha: 1},
	"dragonborn":         {abilityStr: 2, abilityCha: 1},
	"gnome":              {abilityInt: 2},
	"forest gnome":       {abilityInt: 2, abilityDex: 1},
	"rock gnome":         {abilityInt: 2, abilityCon: 1},
	"half-elf":           {abilityCha: 2, "choice1": 1, "choice2": 1},
	"half elf":           {abilityCha: 2, "choice1": 1, "choice2": 1},
	"half-orc":           {abilityStr: 2, abilityCon: 1},
	"half orc":           {abilityStr: 2, abilityCon: 1},
	"tiefling":           {abilityCha: 2, abilityInt: 1},
}

// ClassSavingThrowProficiencies maps class names to their saving throw proficiencies
var ClassSavingThrowProficiencies = map[string][]string{
	"barbarian": {abilityStr, abilityCon},
	"bard":      {abilityDex, abilityCha},
	"cleric":    {abilityWis, abilityCha},
	"druid":     {abilityInt, abilityWis},
	"fighter":   {abilityStr, abilityCon},
	"monk":      {abilityStr, abilityDex},
	"paladin":   {abilityWis, abilityCha},
	"ranger":    {abilityStr, abilityDex},
	"rogue":     {abilityDex, abilityInt},
	"sorcerer":  {abilityCon, abilityCha},
	"warlock":   {abilityWis, abilityCha},
	classWizard: {abilityInt, abilityWis},
}

// SkillToAbility maps skill names to their governing ability
var SkillToAbility = map[string]string{
	"acrobatics":      abilityDex,
	"animal handling": abilityWis,
	"arcana":          abilityInt,
	"athletics":       abilityStr,
	"deception":       abilityCha,
	"history":         abilityInt,
	"insight":         abilityWis,
	"intimidation":    abilityCha,
	"investigation":   abilityInt,
	"medicine":        abilityWis,
	"nature":          abilityInt,
	"perception":      abilityWis,
	"performance":     abilityCha,
	"persuasion":      abilityCha,
	"religion":        abilityInt,
	"sleight of hand": abilityDex,
	"stealth":         abilityDex,
	"survival":        abilityWis,
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
