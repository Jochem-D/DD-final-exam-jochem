package cmd

import (
	"fmt"
	"os"
	"strings"
)

func ExecuteList(args []string) {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Could not read directory:", err)
		os.Exit(1)
	}
	fmt.Println("Characters:")
	found := false
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			fmt.Println(" -", file.Name()[:len(file.Name())-5])
			found = true
		}
	}
	if !found {
		fmt.Println("  (none found)")
	}
}
