package cfzap

import (
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	const TEST_FILE_PATH = "test_config_file"

	// test read config file with wrong type.
	// using cfzap.abc as the config file.
	if _, err := readConfigFile(NewConfigOption(WithFileExt("abc"))); err == nil || err.Error() != "unsupported Config type [abc]" {
		t.Errorf("not expected error: %s", err)
	}

	// test read config file with default name and default type.
	// using cfzap.json as the config file.
	if _, err := readConfigFile(NewConfigOption(WithFileName(""), WithFileExt(""), WithFilePaths(TEST_FILE_PATH))); err != nil {
		t.Errorf("fail to read config file in default type: %s", err)
	}

	// test read config file with default name and json type.
	// using cfzap.json as the config file.
	if _, err := readConfigFile(NewConfigOption(WithFileName(""), WithFileExt("json"), WithFilePaths(TEST_FILE_PATH))); err != nil {
		t.Errorf("fail to read config file in json type: %s", err)
	}

	// test read config file with given name and yaml type.
	// using configfile.yaml as the config file.
	if _, err := readConfigFile(NewConfigOption(WithFileName("configfile"), WithFileExt("yaml"), WithFilePaths(TEST_FILE_PATH))); err != nil {
		t.Errorf("fail to read config file in yaml type: %s", err)
	}
}
