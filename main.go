package main

import (
	"path/filepath"
	"regexp"
	"runtime"

	log "github.com/sirupsen/logrus"

	flag "github.com/spf13/pflag"
)

// TODO: Refactor into generic functions
// TODO: Write tests
// TODO: Validate command: if zettel file is modified and validate all links if no link is broken or head is not consistent with filename
// TODO: Validate command: Check if zettel file deleted (REMOVE) and mark [[links to file]] as bad + report them
// TODO: Create tags index with each tag pointing to files with this tag

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var (
	zettelIDFilenameRegex, _ = regexp.Compile("[0-9]{12}.*\\.md$")
	markdownHeadPrefix, _    = regexp.Compile("^#+\\s+")
)

func main() {

	if runtime.GOOS != "darwin" {
		panic("Because of platform specifics, only MacOS is supported")
	}

	wDir := flag.String("path", ".", "Path to working directory")
	flag.Parse()

	zkDir, err := filepath.Abs(*wDir)
	failOnError(err)

	Serve(zkDir)
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
