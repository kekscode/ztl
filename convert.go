package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func addMarkdownHeadToFile(fn string) (fileHead, filePath string) {
	// Save file content to memory
	fileContent, err := os.ReadFile(fn)
	failOnError(err)

	lines := strings.Split(string(fileContent), "\n")

	//head := fmt.Sprintf("# %s", strings.Split(filepath.Base(fn), ".")[0])
	head := fmt.Sprintf("# %s", strings.Split(filepath.Base(fn), filepath.Ext(fn))[0])
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
	path := strings.Split(fn, filepath.Ext(fn))[0]

	// Strip leading filesystem path
	fileName := filepath.Base(path)

	// Get CWD
	cwd := filepath.Dir(fn)
	failOnError(err)

	// Adjust file name according to head
	if head != fmt.Sprintf("# %s", fileName) {

		newFileName := fmt.Sprintf("%s.md", markdownHeadPrefix.ReplaceAllString(filepath.Base(lines[0]), ""))
		newFileName = filepath.Join(cwd, newFileName)

		err = os.Rename(fn, newFileName)
		failOnError(err)
		return head, newFileName
	}

	return head, fileName
}
