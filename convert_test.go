package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_addMarkdownHeadToFile(t *testing.T) {

	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")
	cwd, _ := os.Getwd()

	type args struct {
		fn string
	}

	a := args{fn: "testdata/202104061620 This is a new markdown file.md"}
	os.WriteFile(a.fn, nil, 0644)

	tests := []struct {
		name         string
		args         args
		wantFileHead string
		wantFilePath string
	}{
		// Test cases
		{
			name:         "Test if first line in new file contains correct title",
			args:         a,
			wantFileHead: "# 202104061620 This is a new markdown file",
			wantFilePath: filepath.Join(cwd, "testdata/202104061620 This is a new markdown file.md"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileHead, gotFilePath := addMarkdownHeadToFile(tt.args.fn)
			if gotFileHead != tt.wantFileHead {
				t.Errorf("addMarkdownHeadToFile() gotFileHead = %v, want %v", gotFileHead, tt.wantFileHead)
			}
			if gotFilePath != tt.wantFilePath {
				t.Errorf("addMarkdownHeadToFile() gotFilePath = %v, want %v", gotFilePath, tt.wantFilePath)
			}
		})
	}
}

func Test_syncFileNameByMarkdownHead(t *testing.T) {

	os.Mkdir("testdata", 0755)
	defer os.RemoveAll("testdata")

	type args struct {
		fn string
	}

	a := args{fn: "testdata/202104061620 This is a new markdown file.md"}
	os.WriteFile(a.fn, []byte("# 202104061620 This is a renamed markdown head"), 0644)

	tests := []struct {
		name         string
		args         args
		wantFileHead string
		wantFilePath string
	}{
		// Test cases
		{
			name:         "Test if file gets renamed according to markdown header",
			args:         a,
			wantFileHead: "# 202104061620 This is a renamed markdown head",
			wantFilePath: filepath.Join(filepath.Base("."), "testdata/202104061620 This is a renamed markdown head.md"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileHead, gotFilePath := syncFileNameByMarkdownHead(tt.args.fn)
			if gotFileHead != tt.wantFileHead {
				t.Errorf("syncFileNameByMarkdownHead() gotFileHead = %v, want %v", gotFileHead, tt.wantFileHead)
			}
			if gotFilePath != tt.wantFilePath {
				t.Errorf("syncFileNameByMarkdownHead() gotFilePath = %v, want %v", gotFilePath, tt.wantFilePath)
			}
		})
	}
}
