package cfzap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigOption(t *testing.T) {
	option := NewConfigOption()

	assert.Equalf(t, ConfigFileName, option.FileName, "default file name should be %q", ConfigFileName)

	option = NewConfigOption(
		WithCreateNew(true),
		WithFileName("test"),
		WithFileExt("json"),
		WithFilePaths("path1", "path2"))

	assert.True(t, option.CreateNew, "createNew should be true")
	assert.Equal(t, "test", option.FileName, "file name should be \"test\"")
	assert.Equal(t, "json", option.FileExt, "file extension should be \"json\"")
	assert.Equal(t, 2, len(option.FilePaths), "there should be 2 paths")
	assert.Equal(t, "path1", option.FilePaths[0], "FilePaths[0] should be \"path1\"")
	assert.Equal(t, "path2", option.FilePaths[1], "FilePaths[1] should be \"path2\"")
}

func TestCompareConfigOption(t *testing.T) {
	option := NewConfigOption()
	assert.True(t, option.equal(option), "same object should be equal.")
	assert.False(t, option.equal(nil), "nil object should not be equal.")

	other := NewConfigOption(WithCreateNew(true))
	assert.False(t, option.equal(other), "property 'CreateNew' is different.")

	other.CreateNew = false
	other.FileExt = "ini"
	assert.False(t, option.equal(other), "property 'FileExt' is different.")
}