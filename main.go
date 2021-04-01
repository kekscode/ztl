package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO:
// validate command: if zettel file is modified and validate all links if no link is broken
// validate command: Check if zettel file deleted (REMOVE) and mark [[links to file]] as bad + report them
func main() {

	if runtime.GOOS != "darwin" {
		panic("Because of platform specifics, only MacOS is supported")
	}

	workDir := flag.String("path", ".", "Path to working directory")
	flag.Parse()

	zettelIDFilenameRegex, _ := regexp.Compile("[0-9]{12}.*\\.md")
	markdownHeadPrefix, _ := regexp.Compile("^#+\\s+")

	watcher, err := fsnotify.NewWatcher()
	failOnError(err)
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// General logging of all modified files:
				//log.Printf("event.Op: %s, event.Name: %s", event.Op, event.Name)

				if zettelIDFilenameRegex.MatchString(event.Name) && event.Op&fsnotify.Create == fsnotify.Create { // CREATE
					log.Printf("Hello from Create: %s", event.Op)
					// Use filename as heading (but strip any extension like .md or else)
					head := fmt.Sprintf("# %s", strings.Split(event.Name, ".")[0])

					// Save file content to memory
					log.Printf("Reading file in 62")
					file, err := os.ReadFile(event.Name)
					failOnError(err)

					lines := strings.Split(string(file), "\n")

					lines[0] = head

					output := strings.Join(lines, "\n")
					log.Printf("Writing file in 71")
					err = os.WriteFile(event.Name, []byte(output), 0644)
					failOnError(err)

					abs, err := filepath.Abs(event.Name)
					failOnError(err)
					log.Printf("Added markdown head \"%s\" to file %s", head, abs)
				}

				if zettelIDFilenameRegex.MatchString(event.Name) && event.Op&fsnotify.Write == fsnotify.Write { // WRITE
					// Check if zettel file first line is modified and correct filename
					log.Printf("Hello from Write: %s", event.Op)

					// Save file content to memory
					log.Printf("Reading file in 121")
					file, err := os.ReadFile(event.Name) // SOMETIMES EMPTY??
					failOnError(err)

					lines := strings.Split(string(file), "\n")
					head := lines[0]

					// Get file name without extension
					fileName := strings.Split(event.Name, ".")[0]

					// Filename and first line of markdown are the same
					if head == fmt.Sprintf("# %s", fileName) {
						log.Printf("OK: Filename \"%s\" and markdown head \"%s\" are in sync", fileName, head)
					}

					// Adjust filename according to heading
					if head != fmt.Sprintf("# %s", fileName) {

						log.Printf("NOK: Filename \"%s\" and markdown head \"%s\" are not in sync", fileName, head)

						log.Printf("Opening file in 141")
						file, err := os.ReadFile(event.Name)
						failOnError(err)
						watcher.Remove(event.Name)

						lines := strings.Split(string(file), "\n")

						newFileName := fmt.Sprintf("%s.md", markdownHeadPrefix.ReplaceAllString(lines[0], ""))
						log.Printf("Renaming file in 149")
						err = os.Rename(event.Name, newFileName)
						watcher.Add(newFileName)
						failOnError(err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*workDir)
	if err != nil {
		log.Fatal(err)
	}

	<-done

}
