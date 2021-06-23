package cfzap

import (
	"testing"
)

const test_config_path = "test_config_file"

func TestLoadAppendersOk(t *testing.T) {
	config, err := readConfigFile("appender_config_ok", "yaml", test_config_path)

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
	config, err := readConfigFile("appender_config_missing", "yaml", test_config_path)

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
