package cfzap

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// last config file information
var (
	lastConfigFileWithoutExt string
	lastConfigFileExit       string
	lastConfigPaths          []string
)

const defaultFilename = "cfzap"

// SetConfigFile sets config file information for logger.
// If don't call this function before calling GetLogger(), it will search current executable path for config file named
// 'cfzap', the extension is applied one by one below.
//
// The supported extensions are:
// "json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "dotenv", "env" and "ini".
// If 'configFileExt' is empty string, the extensions above will be added to the 'configFileWithoutExit' and searched
// one by one.
// It's better don't specify extension for 'configFileWithoutExt'.
// It's better giving a non empty value to 'configFileExt' if 'configFileWithoutExt' has extension.
// Default value "cfzap" will be used if 'configFileWithoutExt' is a empty string.
// If 'configPaths' is empty, it will check current executable path.
//goland:noinspection GoUnusedExportedFunction
func SetConfigFile(configFileWithoutExt string, configFileExt string, configPaths ...string) error {
	configType, typeErr := checkConfigType(configFileExt)
	if typeErr != nil {
		return typeErr
	}

	if s := strings.TrimSpace(configFileWithoutExt); s == "" {
		// using default file name when given value is empty.
		lastConfigFileWithoutExt = defaultFilename
	} else {
		lastConfigFileWithoutExt = s
	}

	lastConfigFileExit = configType
	lastConfigPaths = configPaths

	return nil
}

// readConfigFile reads configuration from specified config file.
// configFileExt must be lowercase or empty string.
// It returns config object and nil when success, otherwise nil and error object.
//func readConfigFile(configFileWithoutExt string, configFileExt string, configPaths ...string) (*viper.Viper, error) {
func readConfigFile(configOption *ConfigOption) (*viper.Viper, error) {
	if configOption == nil {
		configOption = NewConfigOption()
	}

	configType, typeErr := checkConfigType(configOption.FileExt)
	if typeErr != nil {
		return nil, typeErr
	}

	config := viper.New()
	config.SetConfigType(configType)

	if s := strings.TrimSpace(configOption.FileName); s == "" {
		// using default file name when given value is empty.
		config.SetConfigName(defaultFilename)
	} else {
		config.SetConfigName(s)
	}

	if len(configOption.FilePaths) > 0 {
		for _, p := range configOption.FilePaths {
			config.AddConfigPath(strings.TrimSpace(p))
		}
	} else { // using current path when configPaths is empty.
		config.AddConfigPath(".")
	}

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

// checkConfigType checks provided file extension - configFileExt.
// configFileExt is case sensitive or is an empty string - using default value,
// currently is 'json'.
// it returns config type or error.
func checkConfigType(configFileExt string) (string, error) {
	configType := strings.TrimSpace(configFileExt)

	if configType != "" && !stringInSlice(configType, viper.SupportedExts) {
		return "", fmt.Errorf("unsupported Config type [%s]", configFileExt)
	} else {
		return configType, nil
	}
}

func stringInSlice(value string, list []string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}

	return false
}
