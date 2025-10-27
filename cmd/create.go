package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"ddsheetfinal/helpers"
)

func ExecuteCreate(args []string) {
	opts, err := parseCreateFlags(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Build skills
	skillList, err := buildSkillList(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Apply racial bonuses
	strScore, dexScore, conScore, intScore, wisScore, chaScore := applyRacialBonuses(opts)

	// Sort skills and create character
	sort.Strings(skillList)
	char := helpers.Character{
		Name:               opts.Name,
		Race:               opts.Race,
		Class:              opts.Class,
		Background:         opts.Background,
		Level:              opts.Level,
		Str:                strScore,
		Dex:                dexScore,
		Con:                conScore,
		Int:                intScore,
		Wis:                wisScore,
		Cha:                chaScore,
		SkillProficiencies: skillList,
	}

	if err := helpers.SaveCharacter(char); err != nil {
		fmt.Println("Error saving character:", err)
		os.Exit(1)
	}
	fmt.Printf("saved character %s\n", char.Name)
}

type createOptions struct {
	Name       string
	Race       string
	Class      string
	Background string
	Level      int
	Str        int
	Dex        int
	Con        int
	Int        int
	Wis        int
	Cha        int
	Skills     string
}

func parseCreateFlags(args []string) (*createOptions, error) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	// Intentionally empty: we handle flag errors and custom messages ourselves.
	createCmd.Usage = func() {
		// no-op
	}
	name := createCmd.String("name", "", "character name (required)")
	race := createCmd.String("race", "", "character race (required)")
	class := createCmd.String("class", "", "character class (required)")
	background := createCmd.String("background", "acolyte", "character background")
	level := createCmd.Int("level", 1, "character level")
	str := createCmd.Int("str", 10, "strength")
	dex := createCmd.Int("dex", 10, "dexterity")
	con := createCmd.Int("con", 10, "constitution")
	intell := createCmd.Int("int", 10, "intelligence")
	wis := createCmd.Int("wis", 10, "wisdom")
	cha := createCmd.Int("cha", 10, "charisma")
	skills := createCmd.String("skills", "", "comma-separated class skill proficiencies (optional, will be auto-assigned if omitted)")

	if err := createCmd.Parse(args); err != nil {
		return nil, fmt.Errorf("error parsing flags")
	}
	if *name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if *race == "" {
		return nil, fmt.Errorf("race is required")
	}
	if *class == "" {
		return nil, fmt.Errorf("class is required")
	}
	return &createOptions{
		Name:       *name,
		Race:       *race,
		Class:      *class,
		Background: *background,
		Level:      *level,
		Str:        *str,
		Dex:        *dex,
		Con:        *con,
		Int:        *intell,
		Wis:        *wis,
		Cha:        *cha,
		Skills:     *skills,
	}, nil
}

func buildSkillList(opts *createOptions) ([]string, error) {
	bgSkills := helpers.BackgroundSkillProficiencies[strings.ToLower(opts.Background)]
	classSkills := helpers.ClassSkillProficiencies[strings.ToLower(opts.Class)]
	classSkillChoices := helpers.ClassSkillChoices[strings.ToLower(opts.Class)]

	skillList := append([]string{}, bgSkills...)

	if opts.Skills != "" {
		chosen, err := validateChosenSkills(opts.Skills, classSkills, classSkillChoices, opts.Class)
		if err != nil {
			return nil, err
		}
		skillList = append(skillList, chosen...)
		return skillList, nil
	}

	count := 0
	for _, s := range classSkills {
		if count >= classSkillChoices {
			break
		}
		skillList = append(skillList, s)
		count++
	}
	return skillList, nil
}

func validateChosenSkills(raw string, classSkills []string, required int, className string) ([]string, error) {
	chosen := []string{}
	for _, s := range strings.Split(raw, ",") {
		skill := strings.TrimSpace(s)
		ok := false
		for _, cs := range classSkills {
			if strings.EqualFold(skill, cs) {
				ok = true
				break
			}
		}
		if !ok {
			return nil, fmt.Errorf("'%s' is not a valid class skill for %s.", skill, className)
		}
		chosen = append(chosen, skill)
	}
	if len(chosen) != required {
		return nil, fmt.Errorf("You must choose %d class skills.", required)
	}
	return chosen, nil
}

func applyRacialBonuses(opts *createOptions) (int, int, int, int, int, int) {
	raceBonuses := helpers.RacialAbilityBonuses[strings.ToLower(opts.Race)]
	strScore := opts.Str
	if bonus, ok := raceBonuses["str"]; ok {
		strScore += bonus
	}
	dexScore := opts.Dex
	if bonus, ok := raceBonuses["dex"]; ok {
		dexScore += bonus
	}
	conScore := opts.Con
	if bonus, ok := raceBonuses["con"]; ok {
		conScore += bonus
	}
	intScore := opts.Int
	if bonus, ok := raceBonuses["int"]; ok {
		intScore += bonus
	}
	wisScore := opts.Wis
	if bonus, ok := raceBonuses["wis"]; ok {
		wisScore += bonus
	}
	chaScore := opts.Cha
	if bonus, ok := raceBonuses["cha"]; ok {
		chaScore += bonus
	}
	return strScore, dexScore, conScore, intScore, wisScore, chaScore
}
