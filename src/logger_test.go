package cfzap

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	logger, err := GetLogger(nil)

	if err != nil {
		t.Errorf("fail to get new logger: %v", err)
	}

	logger.Debug("this is a test")

}
