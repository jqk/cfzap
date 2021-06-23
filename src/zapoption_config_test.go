package cfzap

import (
	"strconv"
	"testing"
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

	if config, err := readConfigFile(option); err != nil {
		t.Errorf("fail to read config file %d: %s", targetCount, err)
	} else {
		options := loadLogOptions(config)
		count := len(options)

		if count != targetCount {
			t.Errorf("there should be %d options in the list, but got only %d", targetCount, count)
		}
	}
}
