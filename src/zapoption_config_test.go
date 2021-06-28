package cfzap

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadLogOptions(t *testing.T) {
	for i := 0; i < 4; i++ {
		testCase(t, i)
	}
}

func testCase(t *testing.T, targetCount int) {
	filename := "option_config_" + strconv.Itoa(targetCount)
	option := NewConfigOption(
		WithFileName(filename),
		WithFileExt("yaml"),
		WithFilePaths("test_config_file"))

	config, err := readConfigFile(option)
	assert.Nilf(t, err, "fail to read config file for target %d", targetCount)

	options := loadLogOptions(config)
	count := len(options)
	assert.Equalf(t, targetCount, count, "there should be %d options in the list, but got only %d", targetCount, count)
}
