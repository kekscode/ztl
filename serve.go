package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

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
