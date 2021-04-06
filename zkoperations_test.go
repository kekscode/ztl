package main

import "testing"

func Test_addMarkdownHeadToFile(t *testing.T) {
	type args struct {
		fn string
	}
	tests := []struct {
		name         string
		args         args
		wantFileHead string
		wantFilePath string
	}{
		// TODO: Add test cases.
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
	type args struct {
		fn string
	}
	tests := []struct {
		name         string
		args         args
		wantFileHead string
		wantFilePath string
	}{
		// TODO: Add test cases.
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
