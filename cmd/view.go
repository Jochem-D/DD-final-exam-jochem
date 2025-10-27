package cmd

import (
	"ddsheetfinal/helpers"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

const preparedSpellsHeader = "Prepared spells:"
const spellSlotsHeader = "Spell slots:"
const levelLineFmt = "  Level %d: %d\n"

func abilityMod(score int) int { return helpers.AbilityMod(score) }

func ExecuteView(args []string) {
	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	name := viewCmd.String("name", "", "character name (required)")
	enrich := viewCmd.Bool("enrich", false, "fetch extra details from 5e API")

	if err := viewCmd.Parse(args); err != nil || *name == "" {
		fmt.Println("name is required")
		viewCmd.Usage()
		os.Exit(2)
	}

	char := loadCharacter(*name)
	printBasicInfo(char)
	printAbilityScores(char)
	printSkillProficiencies(char)
	printSpellcastingInfo(char, *enrich)
	printEquipment(char)
	printDerivedStats(char, *enrich)
	printPreparedSpells(char)
}

func loadCharacter(name string) helpers.Character {
	filename := helpers.CharacterPath(name)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("character \"%s\" not found\n", name)
		os.Exit(1)
	}

	var char helpers.Character
	if err := json.Unmarshal(data, &char); err != nil {
		fmt.Printf("Could not parse character data: %v\n", err)
		os.Exit(1)
	}
	return char
}

func printBasicInfo(char helpers.Character) {
	fmt.Printf("Name: %s\nClass: %s\nRace: %s\nBackground: %s\nLevel: %d\n",
		char.Name,
		strings.ToLower(char.Class),
		strings.ToLower(char.Race),
		strings.ToLower(char.Background),
		char.Level)
}

func printAbilityScores(char helpers.Character) {
	fmt.Println("Ability scores:")
	fmt.Printf("  STR: %d (%+d)\n", char.Str, abilityMod(char.Str))
	fmt.Printf("  DEX: %d (%+d)\n", char.Dex, abilityMod(char.Dex))
	fmt.Printf("  CON: %d (%+d)\n", char.Con, abilityMod(char.Con))
	fmt.Printf("  INT: %d (%+d)\n", char.Int, abilityMod(char.Int))
	fmt.Printf("  WIS: %d (%+d)\n", char.Wis, abilityMod(char.Wis))
	fmt.Printf("  CHA: %d (%+d)\n", char.Cha, abilityMod(char.Cha))
	fmt.Printf("Proficiency bonus: +%d\n", helpers.GetProficiencyBonus(char.Level))
}

func printSkillProficiencies(char helpers.Character) {
	if len(char.SkillProficiencies) > 0 {
		out := make([]string, 0, len(char.SkillProficiencies))
		for _, s := range char.SkillProficiencies {
			key := strings.ToLower(strings.TrimSpace(s))
			if key == "" {
				continue
			}
			out = append(out, key)
		}
		if len(out) > 0 {
			fmt.Printf("Skill proficiencies: %s\n", strings.Join(out, ", "))
		}
	}
}

func printSpellcastingInfo(char helpers.Character, enrich bool) {
	classLower := strings.ToLower(char.Class)
	fullCasterShow := enrich || (char.Level > 1)

	if fullCasterShow && isFullCaster(classLower) {
		printSpellSlots(&char)
		printSpellcastingAbility(char, classLower)
	} else {
		if classLower == "paladin" || classLower == "warlock" {
			printSpellSlots(&char)
		}
	}
}

func printSpellcastingAbility(char helpers.Character, classLower string) {
	if abbr, ok := spellcastingAbilityForClass(classLower); ok {
		mod := getAbilityModForSpellcasting(char, abbr)
		pb := helpers.GetProficiencyBonus(char.Level)
		fmt.Printf("Spellcasting ability: %s\n", abilityFullName(abbr))
		fmt.Printf("Spell save DC: %d\n", 8+pb+mod)
		fmt.Printf("Spell attack bonus: %+d\n", pb+mod)
	}
}

func getAbilityModForSpellcasting(char helpers.Character, abbr string) int {
	switch abbr {
	case "cha":
		return abilityMod(char.Cha)
	case "wis":
		return abilityMod(char.Wis)
	case "int":
		return abilityMod(char.Int)
	default:
		return 0
	}
}

func printEquipment(char helpers.Character) {
	if char.Weapon != "" {
		fmt.Printf("Main hand: %s\n", char.Weapon)
	}
	if char.OffHand != "" {
		fmt.Printf("Off hand: %s\n", char.OffHand)
	}
	if char.Armor != "" {
		fmt.Printf("Armor: %s\n", char.Armor)
	}
	if char.Shield != "" {
		fmt.Printf("Shield: %s\n", char.Shield)
	}
	if len(char.Inventory) > 0 {
		fmt.Printf("Inventory: %s\n", strings.Join(char.Inventory, ", "))
	}
}

func printDerivedStats(char helpers.Character, enrich bool) {
	classLower := strings.ToLower(char.Class)
	unarmoredBarb := (classLower == "barbarian" || classLower == "monk") && strings.TrimSpace(char.Armor) == ""
	isHighLevelFullCaster := isFullCaster(classLower) && char.Level > 1
	showDerived := enrich || char.Armor != "" || char.Shield != "" || unarmoredBarb || isHighLevelFullCaster

	if showDerived {
		ac, _ := helpers.ComputeAC(char)
		fmt.Printf("Armor class: %d\n", ac)
		fmt.Printf("Initiative bonus: %d\n", helpers.AbilityMod(char.Dex))
		printPassivePerception(char)
	}
}

func printPassivePerception(char helpers.Character) {
	passive := 10 + helpers.AbilityMod(char.Wis)
	for _, s := range char.SkillProficiencies {
		if strings.EqualFold(s, "perception") {
			passive += helpers.GetProficiencyBonus(char.Level)
			break
		}
	}
	fmt.Printf("Passive perception: %d\n", passive)
}

func printPreparedSpells(char helpers.Character) {
	if len(char.PreparedSpells) > 0 {
		fmt.Println(preparedSpellsHeader)
		for _, s := range char.PreparedSpells {
			fmt.Printf("  - %s\n", s)
		}
	}
}

func printSpellSlots(char *helpers.Character) {
	classLower := strings.ToLower(char.Class)

	if classLower == "wizard" || classLower == "cleric" {
		printFullCasterSlots(char)
		return
	}
	if classLower == "warlock" {
		printWarlockSlots(char)
		return
	}
	if classLower == "paladin" || classLower == "ranger" {
		printHalfCasterSlots(char)
		return
	}
}

func printFullCasterSlots(char *helpers.Character) {
	// cantrips known: 3 (L1), 4 (L4+), 5 (L10+)
	cantrips := 3
	if char.Level >= 4 {
		cantrips = 4
	}
	if char.Level >= 10 {
		cantrips = 5
	}
	fullCaster := map[int][]int{
		1:  {2},
		2:  {3},
		3:  {4, 2},
		4:  {4, 3},
		5:  {4, 3, 2},
		6:  {4, 3, 3},
		7:  {4, 3, 3, 1},
		8:  {4, 3, 3, 2},
		9:  {4, 3, 3, 3, 1},
		10: {4, 3, 3, 3, 2},
		11: {4, 3, 3, 3, 2, 1},
		12: {4, 3, 3, 3, 2, 1},
		13: {4, 3, 3, 3, 2, 1, 1},
		14: {4, 3, 3, 3, 2, 1, 1},
		15: {4, 3, 3, 3, 2, 1, 1, 1},
		16: {4, 3, 3, 3, 2, 1, 1, 1},
		17: {4, 3, 3, 3, 2, 1, 1, 1, 1},
		18: {4, 3, 3, 3, 3, 1, 1, 1, 1},
		19: {4, 3, 3, 3, 3, 2, 1, 1, 1},
		20: {4, 3, 3, 3, 3, 2, 2, 1, 1},
	}
	fmt.Println(spellSlotsHeader)
	fmt.Printf("  Level 0: %d\n", cantrips)
	slots := fullCaster[char.Level]
	for lvl := 1; lvl <= len(slots); lvl++ {
		fmt.Printf(levelLineFmt, lvl, slots[lvl-1])
	}
}

func printWarlockSlots(char *helpers.Character) {
	cantrips := 2
	if char.Level >= 4 {
		cantrips = 3
	}
	if char.Level >= 10 {
		cantrips = 4
	}
	pactLevel := map[int]int{
		1: 1, 2: 1,
		3: 2, 4: 2,
		5: 3, 6: 3,
		7: 4, 8: 4,
		9: 5, 10: 5, 11: 5, 12: 5, 13: 5, 14: 5, 15: 5, 16: 5, 17: 5, 18: 5, 19: 5, 20: 5,
	}
	pactSlots := map[int]int{
		1: 1,
		2: 2, 3: 2, 4: 2, 5: 2, 6: 2, 7: 2, 8: 2, 9: 2, 10: 2,
		11: 3, 12: 3, 13: 3, 14: 3, 15: 3, 16: 3,
		17: 4, 18: 4, 19: 4, 20: 4,
	}
	fmt.Println(spellSlotsHeader)
	fmt.Printf("  Level 0: %d\n", cantrips)
	fmt.Printf(levelLineFmt, pactLevel[char.Level], pactSlots[char.Level])
}

func printHalfCasterSlots(char *helpers.Character) {
	halfCaster := map[int][]int{
		1:  {},
		2:  {2},
		3:  {3},
		4:  {3},
		5:  {4, 2},
		6:  {4, 2},
		7:  {4, 3},
		8:  {4, 3},
		9:  {4, 3, 2},
		10: {4, 3, 2},
		11: {4, 3, 3},
		12: {4, 3, 3},
		13: {4, 3, 3, 1},
		14: {4, 3, 3, 1},
		15: {4, 3, 3, 2},
		16: {4, 3, 3, 2},
		17: {4, 3, 3, 3, 1},
		18: {4, 3, 3, 3, 1},
		19: {4, 3, 3, 3, 2},
		20: {4, 3, 3, 3, 2},
	}
	slots := halfCaster[char.Level]
	total := 0
	for _, v := range slots {
		total += v
	}
	if total == 0 {
		return
	}
	fmt.Println(spellSlotsHeader)
	for lvl := 1; lvl <= 5; lvl++ {
		val := 0
		if lvl-1 < len(slots) {
			val = slots[lvl-1]
		}
		fmt.Printf(levelLineFmt, lvl, val)
	}
}

func isFullCaster(classLower string) bool {
	switch classLower {
	case "wizard", "cleric", "druid", "bard", "sorcerer":
		return true
	default:
		return false
	}
}

func spellcastingAbilityForClass(classLower string) (abbr string, ok bool) {
	switch classLower {
	case "bard", "paladin", "sorcerer", "warlock":
		return "cha", true
	case "cleric", "druid", "ranger":
		return "wis", true
	case "wizard", "artificer":
		return "int", true
	default:
		return "", false
	}
}

func abilityFullName(abbr string) string {
	switch abbr {
	case "cha":
		return "charisma"
	case "wis":
		return "wisdom"
	case "int":
		return "intelligence"
	default:
		return abbr
	}
}

func defaultStr(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

func dexBonusString(info helpers.EquipmentInfo) string {
	if !info.DexAllowed {
		return "none"
	}
	if info.MaxDex == nil {
		return "yes"
	}
	return fmt.Sprintf("yes (max %d)", *info.MaxDex)
}
