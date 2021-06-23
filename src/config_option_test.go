package cfzap

import (
	"testing"
)

func TestNewCreateLoggerOption(t *testing.T) {
	option := NewConfigOption()

	if option.FileName != "cfzap" {
		t.Error("default file name should be 'cfzap'")
	}

	option = NewConfigOption(
		WithCreateNew(true),
		WithFileName("abcde"),
		WithFileExt("json"),
		WithFilePaths("path1", "path2"))

	if !option.CreateNew {
		t.Error("createNew should be true")
	}
	if option.FileName != "abcde" {
		t.Error("file name should be 'abcde'")
	}
	if option.FileExt != "json" {
		t.Error("file extension should be 'json'")
	}
	if len(option.FilePaths) != 2 || option.FilePaths[0] != "path1" || option.FilePaths[1] != "path2" {
		t.Error("file paths should be 'path1' and 'path2'")
	}
}
