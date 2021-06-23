package cfzap

import (
	"testing"
)

func TestLoadAppendersOk(t *testing.T) {
	option := NewConfigOption(
		WithFileName("appender_config_ok"),
		WithFileExt("yaml"),
		WithFilePaths("test_config_file"))
	config, err := readConfigFile(option)

	if err != nil {
		t.Errorf("fail to read appender config: %s" + err.Error())
	}

	appenders, errors, err := loadAppenders(config)

	if err != nil {
		t.Errorf("fail to load appender config: %s" + err.Error())
	}

	if len(appenders) != 2 || len(errors) != 0 {
		t.Errorf("appender count or error count is not correct.")
	}
}

func TestLoadAppendersMissing(t *testing.T) {
	option := NewConfigOption(
		WithFileName("appender_config_missing"),
		WithFileExt("yaml"),
		WithFilePaths("test_config_file"))
	config, err := readConfigFile(option)

	if err != nil {
		t.Errorf("fail to read appender config: %s" + err.Error())
	}

	appenders, errors, err := loadAppenders(config)

	if err != nil {
		t.Errorf("fail to load appender config: %s" + err.Error())
	}

	if len(appenders) != 1 || len(errors) != 1 {
		t.Errorf("appender count or error count is not correct.")
	}
}
