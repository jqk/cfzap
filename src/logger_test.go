package cfzap

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	logger, err := GetLogger(true)

	if err != nil {
		t.Errorf("fail to get new logger: %v", err)
	}

	logger.Debug("this is a test")

}
