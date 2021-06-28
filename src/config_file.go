package cfzap

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// readConfigFile reads configuration from specified config file.
// It returns config object and nil when success, otherwise nil and error object.
func readConfigFile(configOption *ConfigOption) (*viper.Viper, error) {
	configType, typeErr := checkConfigType(configOption.FileExt)
	if typeErr != nil {
		return nil, typeErr
	}

	config := viper.New()
	config.SetConfigType(configType)

	if s := strings.TrimSpace(configOption.FileName); s == "" {
		// using default file name when given value is empty.
		config.SetConfigName("cfzap")
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

// checkConfigType checks provided file extension.
// fileExt is case sensitive.
// it returns config type and error object.
func checkConfigType(fileExt string) (string, error) {
	configType := strings.TrimSpace(fileExt)

	// when configType is empty string, it will search Viper supported extensions.
	if configType != "" && !StringInArray(configType, viper.SupportedExts) {
		return "", fmt.Errorf("unsupported Config type [%s]", fileExt)
	} else {
		return configType, nil
	}
}
