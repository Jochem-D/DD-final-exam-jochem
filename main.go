package main

import (
	"fmt"
	"net/http"
	"os"

	"ddsheetfinal/cmd"
	"ddsheetfinal/helpers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		cmd.ExecuteCreate(os.Args[2:])
	case "view":
		cmd.ExecuteView(os.Args[2:])
	case "list":
		cmd.ExecuteList(os.Args[2:])
	case "delete":
		cmd.ExecuteDelete(os.Args[2:])
	case "equip":
		cmd.ExecuteEquip(os.Args[2:])
	case "unequip":
		cmd.ExecuteUnequip(os.Args[2:])
	case "prepare-spell":
		cmd.ExecutePrepareSpell(os.Args[2:])
	case "spell": // alias
		cmd.ExecutePrepareSpell(os.Args[2:])
	case "learn-spell":
		cmd.ExecuteLearnSpell(os.Args[2:])
	case "learnable-spells":
		cmd.ExecuteLearnableSpells(os.Args[2:])
	case "help":
		cmd.ExecuteHelp(os.Args[2:])
	case "serve":
		cmd.ExecuteServe(os.Args[2:])
	case "enrich":
		cmd.ExecuteEnrich(os.Args[2:])
	default:
		fmt.Println("Unknown command")
		os.Exit(2)
	}

	http.HandleFunc("/api/derive", helpers.DeriveHandler)
}
