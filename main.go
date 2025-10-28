package main

import (
	"fmt"
	"os"

	"ddsheetfinal/internal/infrastructure/di"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	// Initialize the dependency injection container
	container := di.NewContainer(
		"data/characters",           // characters directory
		"assets/srd/5e-SRD-Equipment.csv", // equipment CSV
		"assets/srd/5e-SRD-Spells.csv",    // spells CSV
		"data/cache",                // cache directory
	)

	// Route to appropriate handler based on command
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "create":
		container.CreateHandler.Handle(args)
	case "view":
		container.ViewHandler.Handle(args)
	case "list":
		container.ListHandler.Handle(args)
	case "delete":
		container.DeleteHandler.Handle(args)
	case "equip":
		container.EquipHandler.Handle(args)
	case "unequip":
		container.UnequipHandler.Handle(args)
	case "prepare-spell", "spell":
		container.PrepareSpellHandler.Handle(args)
	case "learn-spell":
		container.LearnSpellHandler.Handle(args)
	case "learnable-spells":
		container.LearnableSpellsHandler.Handle(args)
	case "help":
		container.HelpHandler.Handle(args)
	case "serve":
		container.ServeHandler.Handle(args)
	case "enrich":
		container.EnrichHandler.Handle(args)
	default:
		fmt.Println("Unknown command")
		os.Exit(2)
	}
}
