package cfzap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFilePath = "test_config_file"

func TestGetLoggerFromDefaultConfig(t *testing.T) {
	// using default config file: cfzap.json
	logger, err := GetLogger(nil)
	assert.Nil(t, err, "fail to create a new logger from default configuration.")

	// write to log file.
	logger.Debug("create new logger because no config is specified")

	// because we don't provide a new option, it will return exist logger.
	logger2, err := GetLogger(nil)
	assert.Nil(t, err, "fail to get exist logger.")
	assert.Equal(t, logger, logger2, "the returned logger is not the exist one.")

	// write to log file.
	logger.Debug("return exist logger because no config is specified from the second time")

	logger2, err = GetLogger(NewConfigOption(WithCreateNew(true)))
	assert.Nil(t, err, "fail to get create now logger according to ConfigOption.CreateNew.")
	assert.NotEqual(t, logger, logger2, "the new logger should not be same as exist one.")
}

func TestGetLoggerFromDefaultYaml(t *testing.T) {
	// this commented code below reads cfzap.json but nto cfzap.yaml. I think this is a bug in Viper.
	// configOption := NewConfigOption(WithFileName("cfzap"), WithFileExt("yaml"))
	configOption := NewConfigOption(WithFileName("cfzap.yaml"), WithFileExt("yaml"))
	logger, err := GetLogger(configOption)
	assert.Nil(t, err, "fail to create a new logger from cfzap.yaml.")

	// try sugar logger.
	sugarLogger := logger.Sugar()
	sugarLogger.Info("this is a test for SugarLogger")
}

func TestGetLoggerForFileNotFound(t *testing.T) {
	_, err := GetLogger(NewConfigOption(WithFileName("no_file")))
	assert.NotNil(t, err, "should be failed because there's no config file.")
	assert.Equal(t, 0, strings.Index(err.Error(), "Config File \"no_file\" Not Found"), "wrong error message")
}

func TestGetLoggerForMissingAppender(t *testing.T) {
	// don't use NewConfigOption() to create a new instance.
	configOption := &ConfigOption{}
	configOption.FileExt = "yaml"
	configOption.FileName = "appender_config_missing"
	configOption.FilePaths = []string{testFilePath}

	logger, err := GetLogger(configOption)
	assert.Nil(t, err, "should be failed because there's a missing appender.")
	assert.NotNil(t, logger, "there is still a valid appender.")
}

func TestGetLoggerWithoutAppenders(t *testing.T) {
	_, err := GetLogger(NewConfigOption(
		WithFileName("appender_config_no_appenders"),
		WithFileExt("yaml"),
		WithFilePaths(testFilePath)))

	assert.NotNil(t, err, "there's no appenders section defined.")
	assert.Equal(t, "missing section [appenders]", err.Error(), "wrong error message")
}
