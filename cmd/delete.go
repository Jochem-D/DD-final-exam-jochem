package cmd

import (
	"flag"
	"fmt"
	"os"

	"ddsheetfinal/helpers"
)

func ExecuteDelete(args []string) {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	name := deleteCmd.String("name", "", "character name (required)")

	if err := deleteCmd.Parse(args); err != nil || *name == "" {
		fmt.Println("name is required")
		deleteCmd.Usage()
		os.Exit(2)
	}

	filename := helpers.CharacterPath(*name)
	if err := os.Remove(filename); err != nil {
		fmt.Printf("Could not delete character file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("deleted %s\n", *name)
}
