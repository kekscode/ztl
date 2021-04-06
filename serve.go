package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

func addMarkdownHeadToFile(fn string) (fileHead, filePath string) {
	head := fmt.Sprintf("# %s", strings.Split(filepath.Base(fn), ".")[0])

	// Save file content to memory
	fileContent, err := os.ReadFile(fn)
	failOnError(err)

	lines := strings.Split(string(fileContent), "\n")
	lines[0] = head

	output := strings.Join(lines, "\n")

	err = os.WriteFile(fn, []byte(output), 0644)
	failOnError(err)

	abs, err := filepath.Abs(fn)
	failOnError(err)

	return head, abs
}

func syncFileNameByMarkdownHead(fn string) (fileHead, filePath string) {
	// Save file content to memory
	fileContent, err := os.ReadFile(fn)

	failOnError(err)

	lines := strings.Split(string(fileContent), "\n")
	head := lines[0]

	// Strip file extension
	path := strings.Split(fn, ".")[0]

	// Strip leading filesystem path
	fileName := filepath.Base(path)

	// Get CWD
	cwd, err := filepath.Abs(".")
	failOnError(err)

	// Adjust file name according to head
	if head != fmt.Sprintf("# %s", fileName) {

		log.Printf("File name \"%s\" and markdown head \"%s\" are not consistent. Adjusting.", fileName, head)

		newFileName := fmt.Sprintf("%s.md", markdownHeadPrefix.ReplaceAllString(filepath.Base(lines[0]), ""))

		newFileName = filepath.Join(cwd, newFileName)

		err = os.Rename(fn, newFileName)
		failOnError(err)
		return head, newFileName
	}

	// Filename and first line of markdown are the same
	if head == fmt.Sprintf("# %s", fileName) {
		log.Printf("File name \"%s\" and markdown head \"%s\" are consistent.", filepath.Join(cwd, fileName), head)
	}

	return head, fileName
}

func Serve(dir string) {
	watcher, err := fsnotify.NewWatcher()
	failOnError(err)
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		log.Println("Starting ztl server...")
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Use file name w/o extension as markdown head on file creation
				if zettelIDFilenameRegex.MatchString(event.Name) && event.Op&fsnotify.Create == fsnotify.Create { // CREATE
					head, abs := addMarkdownHeadToFile(event.Name)
					log.Printf("Added markdown head \"%s\" to new file %s", head, abs)
				}

				// Check if zettel file first line is modified and sync file name accordingly
				if zettelIDFilenameRegex.MatchString(event.Name) && event.Op&fsnotify.Write == fsnotify.Write { // WRITE
					newHead, newFilePath := syncFileNameByMarkdownHead(event.Name)
					log.Printf("Renamed file \"%s\" according to new markdown head %s", newFilePath, newHead)
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
