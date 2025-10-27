package cmd

import (
	"ddsheetfinal/helpers"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const errSavingCharacter = "Error saving character:"
const flagCharNameReq = "character name (required)"

// ExecuteLearnSpell adds a learned spell to the character if valid for their class.
func ExecuteLearnSpell(args []string) {
	fs := flag.NewFlagSet("learn-spell", flag.ExitOnError)
	name := fs.String("name", "", flagCharNameReq)
	spell := fs.String("spell", "", "spell name (required, quote if it has spaces)")
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if strings.TrimSpace(*name) == "" || strings.TrimSpace(*spell) == "" {
		fmt.Println("name and spell are required")
		os.Exit(1)
	}

	filename := helpers.CharacterPath(*name)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read character file: %v\n", err)
		os.Exit(1)
	}

	var c helpers.Character
	if err := json.Unmarshal(data, &c); err != nil {
		fmt.Printf("Could not parse character data: %v\n", err)
		os.Exit(1)
	}

	cls := strings.ToLower(strings.TrimSpace(c.Class))
	if !isCasterClass(cls) {
		fmt.Println("this class can't cast spells")
		os.Exit(1)
	}
	// Prepared casters (cleric, druid, paladin, wizard) don't learn known spells.
	if isPreparedClass(cls) {
		// exact required output
		fmt.Println("this class prepares spells and can't learn them")
		os.Exit(1)
	}

	if !helpers.IsValidSpell(*spell) {
		fmt.Printf("Spell '%s' not found in 5e-SRD-Spells.csv\n", *spell)
		os.Exit(1)
	}
	ok, err := isSpellForClass(*spell, c.Class)
	if err != nil {
		fmt.Println("Error validating spell for class:", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Printf("'%s' is not available to the %s class.\n", *spell, strings.ToLower(c.Class))
		os.Exit(1)
	}

	// Already learned?
	for _, s := range c.Spells {
		if strings.EqualFold(s, *spell) {
			fmt.Printf("%s already knows the spell '%s'.\n", c.Name, *spell)
			return
		}
	}

	c.Spells = append(c.Spells, *spell)
	sort.Strings(c.Spells)

	if err := helpers.SaveCharacter(c); err != nil {
		fmt.Println("Error saving character:", err)
		os.Exit(1)
	}

	fmt.Printf("Learned spell %s\n", *spell)
}

// ExecutePrepareSpell prepares or unprepares a learned spell for the character.
// Examples:
//
//	dndcsg prepare-spell -name Gandalf -spell "burning hands"
//	dndcsg spell -name Gandalf -spell "burning hands"     (alias via main)
func ExecutePrepareSpell(args []string) {
	name, spell, remove, err := parsePrepareFlags(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c, err := loadCharacterFromFile(name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cls := strings.ToLower(strings.TrimSpace(c.Class))
	if !isCasterClass(cls) {
		fmt.Println("this class can't cast spells")
		os.Exit(1)
	}
	if isKnownOnlyClass(cls) {
		fmt.Println("this class learns spells and can't prepare them")
		os.Exit(1)
	}

	target, err := validateSpellForPrepare(c, spell)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := handlePrepareAction(c, cls, remove, target); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parsePrepareFlags(args []string) (string, string, bool, error) {
	fs := flag.NewFlagSet("prepare-spell", flag.ExitOnError)
	name := fs.String("name", "", flagCharNameReq)
	spell := fs.String("spell", "", "spell name (required, quote if it has spaces)")
	remove := fs.Bool("remove", false, "unprepare the spell")
	unprepare := fs.Bool("unprepare", false, "alias for -remove")
	if err := fs.Parse(args); err != nil {
		return "", "", false, err
	}
	if strings.TrimSpace(*name) == "" || strings.TrimSpace(*spell) == "" {
		return "", "", false, fmt.Errorf("name and spell are required")
	}
	if *unprepare {
		*remove = true
	}
	return *name, *spell, *remove, nil
}

func loadCharacterFromFile(name string) (*helpers.Character, error) {
	filename := helpers.CharacterPath(name)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Could not read character file: %v", err)
	}
	var c helpers.Character
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("Could not parse character data: %v", err)
	}
	return &c, nil
}

func validateSpellForPrepare(c *helpers.Character, spell string) (string, error) {
	// ensure the spell exists and is available to the class, and not a cantrip
	if !helpers.IsValidSpell(spell) {
		return "", fmt.Errorf("Spell '%s' not found in 5e-SRD-Spells.csv", spell)
	}
	ok, err := isSpellForClass(spell, c.Class)
	if err != nil {
		return "", fmt.Errorf("Error validating spell for class: %v", err)
	}
	if !ok {
		return "", fmt.Errorf("'%s' is not available to the %s class.", spell, strings.ToLower(c.Class))
	}
	spellLvl, err := getSpellLevel(spell)
	if err != nil {
		return "", err
	}
	cls := strings.ToLower(strings.TrimSpace(c.Class))
	maxLvl := maxSlotLevelForClass(cls, c.Level)
	if spellLvl > maxLvl {
		return "", fmt.Errorf("the spell has higher level than the available spell slots")
	}
	return normalizeSpellName(spell), nil
}

func handlePrepareAction(c *helpers.Character, cls string, remove bool, target string) error {
	if !remove {
		if cls == "wizard" && !hasLearned(c, target) {
			c.Spells = append(c.Spells, target)
			sort.Strings(c.Spells)
		}
		c.PreparedSpells = appendIfMissingPrepared(c.PreparedSpells, target)
		if err := helpers.SaveCharacter(*c); err != nil {
			return fmt.Errorf("%s %v", errSavingCharacter, err)
		}
		fmt.Printf("Prepared spell %s\n", target)
		return nil
	}
	// unprepare
	removed := false
	for i, s := range c.PreparedSpells {
		if strings.EqualFold(s, target) {
			c.PreparedSpells = append(c.PreparedSpells[:i], c.PreparedSpells[i+1:]...)
			removed = true
			break
		}
	}
	if !removed {
		return fmt.Errorf("%s did not have %s prepared", c.Name, target)
	}
	if err := helpers.SaveCharacter(*c); err != nil {
		return fmt.Errorf("%s %v", errSavingCharacter, err)
	}
	fmt.Printf("Unprepared spell %s\n", target)
	return nil
}

// appendIfMissingPrepared appends target to the prepared list unless a case-insensitive match exists.
func appendIfMissingPrepared(prepared []string, target string) []string {
	for _, p := range prepared {
		if strings.EqualFold(p, target) {
			return prepared
		}
	}
	return append(prepared, target)
}

// removePrepared removes case-insensitive matches of target from the prepared list.
func removePrepared(prepared []string, target string) []string {
	out := []string{}
	for _, p := range prepared {
		if strings.EqualFold(p, target) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// Helpers

func normalizeSpellName(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

func hasLearned(c *helpers.Character, normalizedSpell string) bool {
	for _, s := range c.Spells {
		if strings.EqualFold(normalizeSpellName(s), normalizedSpell) {
			return true
		}
	}
	return false
}

func isCasterClass(cls string) bool {
	switch cls {
	case "bard", "cleric", "druid", "paladin", "ranger", "sorcerer", "warlock", "wizard":
		return true
	default:
		return false
	}
}

func isPreparedClass(cls string) bool {
	switch cls {
	case "cleric", "druid", "paladin", "wizard":
		return true
	default:
		return false
	}
}

func isKnownOnlyClass(cls string) bool {
	switch cls {
	case "bard", "sorcerer", "warlock", "ranger":
		return true
	default:
		return false
	}
}

func isSpellForClass(spellName, className string) (bool, error) {
	records, err := readSpellsCSV()
	if err != nil {
		return false, err
	}
	nameIdx, classIdx, err := spellCSVHeaderIndexes(records[0])
	if err != nil {
		return false, err
	}
	want := strings.ToLower(strings.TrimSpace(spellName))
	cls := strings.ToLower(strings.TrimSpace(className))
	for _, r := range records[1:] {
		if len(r) <= classIdx {
			continue
		}
		if strings.ToLower(strings.TrimSpace(r[nameIdx])) != want {
			continue
		}
		classes := strings.ToLower(strings.TrimSpace(r[classIdx]))
		return strings.Contains(classes, cls), nil
	}
	return false, nil
}

func readSpellsCSV() ([][]string, error) {
	f, err := os.Open("5e-SRD-Spells.csv")
	if err != nil {
		return nil, fmt.Errorf("could not open 5e-SRD-Spells.csv: %w", err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil || len(records) == 0 {
		if err == nil {
			err = fmt.Errorf("empty CSV")
		}
		return nil, fmt.Errorf("could not read 5e-SRD-Spells.csv: %w", err)
	}
	return records, nil
}

func spellCSVHeaderIndexes(header []string) (nameIdx, classIdx int, err error) {
	nameIdx, classIdx = -1, -1
	for i, col := range header {
		switch strings.ToLower(strings.TrimSpace(col)) {
		case "name":
			nameIdx = i
		case "class":
			classIdx = i
		}
	}
	if nameIdx == -1 || classIdx == -1 {
		return -1, -1, fmt.Errorf("CSV must have 'name' and 'class' columns")
	}
	return nameIdx, classIdx, nil
}

// getSpellLevel returns the numeric spell level from 5e-SRD-Spells.csv (0..9).
func getSpellLevel(spellName string) (int, error) {
	records, err := readSpellsCSV()
	if err != nil {
		return 0, err
	}
	nameIdx, lvlIdx := -1, -1
	for i, col := range records[0] {
		switch strings.ToLower(strings.TrimSpace(col)) {
		case "name":
			nameIdx = i
		case "level":
			lvlIdx = i
		}
	}
	if nameIdx == -1 || lvlIdx == -1 {
		return 0, fmt.Errorf("CSV must have 'name' and 'level' columns")
	}
	want := strings.ToLower(strings.TrimSpace(spellName))
	for _, r := range records[1:] {
		if len(r) <= lvlIdx {
			continue
		}
		if strings.ToLower(strings.TrimSpace(r[nameIdx])) != want {
			continue
		}
		v, err := strconv.Atoi(strings.TrimSpace(r[lvlIdx]))
		if err != nil {
			return 0, fmt.Errorf("invalid spell level for %q", spellName)
		}
		return v, nil
	}
	return 0, fmt.Errorf("spell %q not found", spellName)
}

// maxSlotLevelForClass returns the highest spell level of slots available to the class at the given level.
func maxSlotLevelForClass(cls string, level int) int {
	fullMax := map[int]int{
		1: 1, 2: 1,
		3: 2, 4: 2,
		5: 3, 6: 3,
		7: 4, 8: 4,
		9: 5, 10: 5,
		11: 6, 12: 6,
		13: 7, 14: 7,
		15: 8, 16: 8,
		17: 9, 18: 9, 19: 9, 20: 9,
	}
	halfMax := map[int]int{
		1: 0, 2: 1, 3: 1, 4: 1,
		5: 2, 6: 2, 7: 2, 8: 2,
		9: 3, 10: 3, 11: 3, 12: 3,
		13: 4, 14: 4, 15: 4, 16: 4,
		17: 5, 18: 5, 19: 5, 20: 5,
	}
	switch cls {
	case "wizard", "cleric", "druid", "bard", "sorcerer":
		return fullMax[level]
	case "paladin", "ranger":
		return halfMax[level]
	case "warlock":
		// not used (warlocks don't prepare), but for completeness:
		pactLevel := map[int]int{
			1: 1, 2: 1, 3: 2, 4: 2, 5: 3, 6: 3, 7: 4, 8: 4, 9: 5,
			10: 5, 11: 5, 12: 5, 13: 5, 14: 5, 15: 5, 16: 5, 17: 5, 18: 5, 19: 5, 20: 5,
		}
		return pactLevel[level]
	default:
		return 0
	}
}

// ExecuteLearnableSpells prints spells the character could learn (for learning classes).
func ExecuteLearnableSpells(args []string) {
	fs := flag.NewFlagSet("learnable-spells", flag.ExitOnError)
	name := fs.String("name", "", flagCharNameReq)
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if strings.TrimSpace(*name) == "" {
		fmt.Println("name is required")
		os.Exit(1)
	}
	printLearnableFor(*name)
}

func printLearnableFor(name string) {
	filename := helpers.CharacterPath(name)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Could not read character file: %v\n", err)
		os.Exit(1)
	}
	var c helpers.Character
	if err := json.Unmarshal(data, &c); err != nil {
		fmt.Printf("Could not parse character data: %v\n", err)
		os.Exit(1)
	}

	records, err := readSpellsCSV()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	learned := map[string]bool{}
	for _, s := range c.Spells {
		learned[strings.ToLower(strings.TrimSpace(s))] = true
	}

	nameIdx, classIdx, err := spellCSVHeaderIndexes(records[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cls := strings.ToLower(strings.TrimSpace(c.Class))
	isPrepared := isPreparedClass(cls)
	isKnownOnly := isKnownOnlyClass(cls)

	fmt.Println("Learnable spells:")
	found := printLearnableSpells(records[1:], nameIdx, classIdx, cls, isPrepared, isKnownOnly, learned)
	if !found {
		fmt.Println("  (none)")
	}
}

// printLearnableSpells prints spells that are learnable for the character and returns true if any are found.
func printLearnableSpells(records [][]string, nameIdx, classIdx int, cls string, isPrepared, isKnownOnly bool, learned map[string]bool) bool {
	found := false
	for _, rec := range records {
		if len(rec) <= classIdx {
			continue
		}
		spellName := strings.TrimSpace(rec[nameIdx])
		classes := strings.ToLower(strings.TrimSpace(rec[classIdx]))

		if !strings.Contains(classes, cls) {
			continue
		}

		learnedSpell := learned[strings.ToLower(spellName)]
		if (isPrepared && cls != "wizard") ||
			(isKnownOnly && !learnedSpell) ||
			(cls == "wizard" && !learnedSpell) {
			fmt.Printf(" - %s\n", spellName)
			found = true
		}
	}
	return found
}
