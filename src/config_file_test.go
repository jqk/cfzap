package cfzap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	var err error

	// test read config file with wrong type.
	// using cfzap.abc as the config file.
	_, err = readConfigFile(NewConfigOption(WithFileExt("abc")))
	assert.NotNil(t, err)
	assert.Equal(t, "unsupported Config type [abc]", err.Error(), "not expected error.")

	// test read config file with default name and default type.
	// using cfzap.json as the config file.
	_, err = readConfigFile(NewConfigOption(WithFileName(""), WithFileExt(""), WithFilePaths(testFilePath)))
	assert.Nil(t, err, "fail to read config file in default type.")

	// test read config file with default name and json type.
	// using cfzap.json as the config file.
	_, err = readConfigFile(NewConfigOption(WithFileName(""), WithFileExt("json"), WithFilePaths(testFilePath)))
	assert.Nil(t, err, "fail to read config file in json type.")

	// test read config file with given name and yaml type.
	// using configfile.yaml as the config file.
	_, err = readConfigFile(NewConfigOption(WithFileName("configfile"), WithFileExt("yaml"), WithFilePaths(testFilePath)))
	assert.Nil(t, err, "fail to read config file in yaml type.")
}
