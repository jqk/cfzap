package cfzap

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// loadLogOptions loads additional option from config file.
// return empty option list when there's no entry.
func loadLogOptions(config *viper.Viper) []zap.Option {
	// 'options' is the fixed top level key.
	section := config.Sub("options")
	options := []zap.Option{}

	if section == nil {
		// return empty option list when there is 'option' section.
		return options
	}

	// 'caller' is the fixed key inside options.
	if section.GetBool("caller") {
		options = append(options, zap.AddCaller())
	}

	// 'development' is the fixed key inside options.
	if section.GetBool("development") {
		options = append(options, zap.Development())
	}

	// 'fields' is the fixed key inside options.
	section = section.Sub("fields")
	if section != nil {
		keys := section.AllKeys()

		if count := len(keys); count > 0 {
			fields := make([]zap.Field, count)
			for i, key := range keys {
				fields[i] = zap.String(key, section.GetString(key))
			}

			options = append(options, zap.Fields(fields...))
		}
	}

	return options
}
