package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func Serve(dir string) {
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

				// Use file name w/o extension as markdown head on file creation
				if zettelIDFilenameRegex.MatchString(event.Name) && event.Op&fsnotify.Create == fsnotify.Create { // CREATE
					head := fmt.Sprintf("# %s", strings.Split(filepath.Base(event.Name), ".")[0])

					// Save file content to memory
					fileContent, err := os.ReadFile(event.Name)
					failOnError(err)

					lines := strings.Split(string(fileContent), "\n")
					lines[0] = head

					output := strings.Join(lines, "\n")

					err = os.WriteFile(event.Name, []byte(output), 0644)
					failOnError(err)

					abs, err := filepath.Abs(event.Name)
					failOnError(err)
					log.Printf("Added markdown head \"%s\" to new file %s", head, abs)
				}

				// Check if zettel file first line is modified and sync file name accordingly
				if zettelIDFilenameRegex.MatchString(event.Name) && event.Op&fsnotify.Write == fsnotify.Write { // WRITE

					// Save file content to memory
					fileContent, err := os.ReadFile(event.Name)
					failOnError(err)

					lines := strings.Split(string(fileContent), "\n")
					head := lines[0]

					// Strip file extension
					filePath := strings.Split(event.Name, ".")[0]

					// Strip leading filesystem path
					fileName := filepath.Base(filePath)

					// Filename and first line of markdown are the same
					if head == fmt.Sprintf("# %s", fileName) {
						log.Printf("File name \"%s\" and markdown head \"%s\" are consistent.", fileName, head)
					}

					// Adjust file name according to head
					if head != fmt.Sprintf("# %s", fileName) {

						log.Printf("File name \"%s\" and markdown head \"%s\" are not consistent. Adjusting.", fileName, head)

						watcher.Remove(event.Name)

						newFile := fmt.Sprintf("%s.md", markdownHeadPrefix.ReplaceAllString(filepath.Base(lines[0]), ""))
						newFile = filepath.Join(dir, newFile)

						err = os.Rename(event.Name, newFile)
						watcher.Add(newFile)
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

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
