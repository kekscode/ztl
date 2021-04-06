package main

import (
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/integrii/flaggy"

	log "github.com/sirupsen/logrus"
)

// FIXME: On MacOS (case-insensitive filesystem) a rename from `202104061547 file.md` to `202104061547 File.md` creates an endless loop in serve.go
//        The same happens when the markdown title is adjusted from `# 202104061547 file` to `# 202104061547 File`
//        A "breaker" is needed here in order to avoid the endless loop.
// TODO: Validate command: if zettel file is modified and validate all links if no link is broken or head is not consistent with filename
// TODO: Validate command: Check if zettel file deleted (REMOVE) and mark [[links to file]] as bad + report them
// TODO: Create tags index with each tag pointing to files with this tag
//       <https://rosettacode.org/wiki/Inverted_index#Go>

const (
	version = "0.0.0-unreleased"
)

var (
	// Zettelkasten-related pattern
	zettelIDFilenameRegex, _ = regexp.Compile("[0-9]{12}.*\\.md$")
	markdownHeadPrefix, _    = regexp.Compile("^#+\\s+")

	// Config related
	workingDirectory = "."
	fixIssues        = false

	// Keep subcommands as globals so you can easily check if they were used later on.
	subcmdServe    *flaggy.Subcommand
	subcmdValidate *flaggy.Subcommand
)

func init() {
	// Only MacOS is tested for the time being
	if runtime.GOOS != "darwin" {
		panic("Because of platform specifics, only MacOS is supported")
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	// CLI
	flaggy.SetName("ztl")
	flaggy.SetDescription("A little server and CLI tool to keep a zettelkasten in shape.")
	flaggy.SetVersion(version)

	subcmdValidate := flaggy.NewSubcommand("validate")
	subcmdValidate.Description = "Validate and (optionally) fix issues with your zettelkasten."
	subcmdValidate.String(&workingDirectory, "w", "work-dir", "Working directory with your zettelkasten files.")
	subcmdValidate.Bool(&fixIssues, "f", "fix-issues", "Validates and fixes issues in one step.")
	flaggy.AttachSubcommand(subcmdValidate, 1)

	subcmdServe := flaggy.NewSubcommand("serve")
	subcmdServe.Description = "Start a server which reacts to file changes in your zettelkasten."
	subcmdServe.String(&workingDirectory, "w", "work-dir", "Working directory with your zettelkasten files.")
	flaggy.AttachSubcommand(subcmdServe, 1)

	flaggy.Parse()

	zkDir, err := filepath.Abs(workingDirectory)
	failOnError(err)

	if subcmdServe.Used {
		Serve(zkDir) // blocks
	}

	flaggy.ShowHelpAndExit("Please provide a subcommand")
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
