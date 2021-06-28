package cfzap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	logger, err := GetLogger(nil)
	assert.Nil(t, err, "fail to get new logger with default configuration.")

	logger.Debug("this is a test 1")

	logger2, err := GetLogger(nil)
	assert.Nil(t, err, "fail to get exist logger.")
	assert.Equal(t, logger, logger2, "the returned logger is not the exist one.")

	logger.Debug("this is a test 2")

	//configOption := NewConfigOption(WithFileName("cfzap"), WithFileExt("yaml"))
	configOption := NewConfigOption(WithFileName("cfzap.yaml"), WithFileExt("yaml"))
	logger, err = GetLogger(configOption)
	assert.Nil(t, err, "fail to get new logger with yaml configuration.")

	sugarLogger := logger.Sugar()
	sugarLogger.Info("this is a test 3")
}
