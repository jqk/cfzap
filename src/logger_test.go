package cfzap

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	logger, err := GetLogger(nil)

	if err != nil {
		t.Errorf("fail to get new logger: %v", err)
	}

	logger.Debug("this is a test 1")

	configOption := NewConfigOption(WithFileExt("yaml"), WithFileName("cfzap2"))
	logger, err = GetLogger(configOption)

	if err != nil {
		t.Errorf("fail to get exist logger: %v", err)
	}

	sugarLogger := logger.Sugar()
	sugarLogger.Info("this is a test 2")
}
