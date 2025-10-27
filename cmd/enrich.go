// cmd/enrich.go
package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"ddsheetfinal/helpers"
)

func ExecuteEnrich(args []string) {
	inplace, force, fetchAll, name := parseEnrichArgs(args)

	baseCtx := context.Background()

	if fetchAll {
		if err := runFetchAll(baseCtx); err != nil {
			fmt.Printf("fetch-all failed: %v\n", err)
			os.Exit(1)
		}
		// If no name was provided, we ran in fetch-only mode: exit successfully.
		if name == "" {
			fmt.Println("fetch-only completed")
			return
		}
	}

	// 60s is reasonable for a single-character enrich run.
	ctx, cancel := context.WithTimeout(baseCtx, 60*time.Second)
	defer cancel()

	out, err := helpers.RunEnrich(ctx, name, inplace, force, &helpers.StdoutLog{})
	if err != nil {
		fmt.Printf("Enrich failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Finished Wrote:", out)
}

// parseEnrichArgs parses flags and the optional name argument from args.
func parseEnrichArgs(args []string) (inplace, force, fetchAll bool, name string) {
	fs := flag.NewFlagSet("enrich", flag.ExitOnError)
	inplacePtr := fs.Bool("inplace", false, "update the original JSON instead of writing to enrichments/<name>.json")
	forcePtr := fs.Bool("force", false, "force fetching/caching even if enriched file exists")
	fetchAllPtr := fs.Bool("fetch-all", false, "fetch all spells and equipment from the API into the .cache before enriching")
	fs.Parse(args)

	// The stdlib flag package stops parsing flags after the first
	// non-flag argument, so callers may pass flags after the name
	// (e.g. `enrich Gandalf --force`). To be forgiving, scan the
	// raw args for known flags and set them if present.
	for _, a := range args {
		if a == "--force" || a == "-force" {
			*forcePtr = true
		}
		if a == "--inplace" || a == "-inplace" {
			*inplacePtr = true
		}
		if a == "--fetch-all" || a == "-fetch-all" {
			*fetchAllPtr = true
		}
	}

	inplace = *inplacePtr
	force = *forcePtr
	fetchAll = *fetchAllPtr

	if fs.NArg() >= 1 {
		name = strings.TrimSpace(fs.Arg(0))
	} else {
		name = ""
	}
	// If neither fetch-all nor a name were provided, show usage and exit.
	if !fetchAll && name == "" {
		fmt.Println("Usage: go run main.go enrich <name> [--inplace] or: go run main.go enrich --fetch-all")
		os.Exit(1)
	}
	return
}

// runFetchAll performs the long-running fetch-all operation under a timeout.
func runFetchAll(baseCtx context.Context) error {
	fmt.Println("fetch-all: fetching all spells and equipment into .cache (this may take a while)")
	ctxFetch, cancelFetch := context.WithTimeout(baseCtx, 10*time.Minute)
	defer cancelFetch()
	if err := helpers.FetchAllReferences(ctxFetch, &helpers.StdoutLog{}); err != nil {
		return err
	}
	fmt.Println("fetch-all: done")
	return nil
}
