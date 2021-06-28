package cfzap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAppendersOk(t *testing.T) {
	option := NewConfigOption(
		WithFileName("appender_config_ok"),
		WithFileExt("yaml"),
		WithFilePaths("test_config_file"))
	config, err := readConfigFile(option)

	assert.Nil(t, err, "fail to read appender config")

	appenders, errors, err := loadAppenders(config)
	assert.Nil(t, err, "fail to load appender config")
	assert.Equal(t, 2, len(appenders), "appender count should be 2")
	assert.Equal(t, 0, len(errors), "error count should be 2")
}

func TestLoadAppendersMissing(t *testing.T) {
	option := NewConfigOption(
		WithFileName("appender_config_missing"),
		WithFileExt("yaml"),
		WithFilePaths("test_config_file"))
	config, err := readConfigFile(option)

	assert.Nil(t, err, "fail to read appender config")

	appenders, errors, err := loadAppenders(config)
	assert.Nil(t, err, "fail to load appender config")
	assert.Equal(t, 1, len(appenders), "appender count should be 1")
	assert.Equal(t, 1, len(errors), "error count should be 1")
}
