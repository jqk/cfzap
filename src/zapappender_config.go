package cfzap

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// appenderConfig includes config info for creating logger.
type appenderConfig struct {
	// zapcore.writeSyncer used for creating logger.
	writeSyncer *zapcore.WriteSyncer
	// zapcore.encoder used for creating logger.
	encoder *zapcore.Encoder
	// the zap.AtomicLevel used for creating logger.
	logLevel zap.AtomicLevel
	// the name of appender, for debug only.
	name string
	// the zapcore.EncoderConfig needed by Encoder.
	encoderConfig *zapcore.EncoderConfig
}

// loadAppenders loads all appenders defined in config file section 'appenders'.
// it returns the successful loaded appender list , failed appender list and error object.
func loadAppenders(config *viper.Viper) (map[string]*appenderConfig, map[string]error, error) {
	// 'appenders' is the fixed top level key and cannot be ignored.
	const SECTION_NAME = "appenders"

	if !config.IsSet(SECTION_NAME) {
		return nil, nil, fmt.Errorf("missing section [%s]", SECTION_NAME)
	}

	appenderNames := config.Get(SECTION_NAME).([]interface{})

	// there is at least one appender.
	if len(appenderNames) == 0 {
		return nil, nil, fmt.Errorf("no appender is defined in section [%s]", SECTION_NAME)
	}

	appenders := make(map[string]*appenderConfig)
	errorAppenders := make(map[string]error)

	// load appenders from config file and put them into a map to filter duplication.
	for _, v := range appenderNames {
		w := v.(string)
		appenders[strings.TrimSpace(w)] = nil
	}

	// load each appender info from config file.
	// if error happened, put it into error appender list.
	for appenderName := range appenders {
		if appender, err := loadAppender(config, appenderName); err == nil {
			appenders[appenderName] = appender
		} else {
			errorAppenders[appenderName] = err
		}
	}

	// remove appenders which is failed to load.
	if len(errorAppenders) > 0 {
		for k := range errorAppenders {
			delete(appenders, k)
		}
	}

	// there should be at least one successful loaded appender.
	if len(appenders) == 0 {
		return nil, errorAppenders, fmt.Errorf("failed to load all %d appenders", len(errorAppenders))
	}

	return appenders, errorAppenders, nil
}

// loadAppender loads each appender according to its name.
// It returns appenderConfig object and error when the entry was missing.
func loadAppender(config *viper.Viper, appenderName string) (*appenderConfig, error) {
	appenderSection := config.Sub(appenderName)
	if appenderSection == nil {
		return nil, fmt.Errorf("cannot find the entry for appender [%s]", appenderName)
	}

	appender := new(appenderConfig)
	appender.name = appenderName

	if err := loadAppenderEncoderConfig(config, appender, appenderSection, appenderName); err != nil {
		return nil, err
	}
	if err := loadAppenderWriteSyncer(config, appender, appenderSection, appenderName); err != nil {
		return nil, err
	}

	loadAppenderLogLevel(appender, appenderSection)
	// must be called after loadAppenderEncoderConfig() because Encoder needs EncoderConfig.
	loadAppenderEncoder(appender, appenderSection)

	return appender, nil
}

// loadAppenderWriteSyncer loads WriteSyncer from config file.
// It returns error when the entry was missing.
func loadAppenderWriteSyncer(config *viper.Viper, appender *appenderConfig, appenderSection *viper.Viper, appenderName string) error {
	s, err := getRequiredString(appenderSection, appenderName, "target")
	if err != nil {
		return err
	}

	var syncer zapcore.WriteSyncer

	if strings.EqualFold(s, "stdout") {
		syncer = zapcore.AddSync(os.Stdout)
	} else if strings.EqualFold(s, "stderr") {
		syncer = zapcore.AddSync(os.Stderr)
	} else {
		section := config.Sub(s)
		if section == nil {
			return fmt.Errorf("the value of [%s.target] is [%s], but the entry was missing", appenderName, s)
		}

		writer, err := loadLumberjack(section)
		if err != nil {
			return err
		}

		syncer = zapcore.AddSync(writer)
	}

	appender.writeSyncer = &syncer

	return nil
}

// loadLumberjack loads lumberjack.Logger as io.Writer from config file.
// It returns error when it failed to create log file path.
func loadLumberjack(section *viper.Viper) (*lumberjack.Logger, error) {
	writer := new(lumberjack.Logger)

	// user must provide valid values.
	writer.Compress = section.GetBool("compress")
	writer.LocalTime = section.GetBool("localTime")
	writer.MaxAge = section.GetInt("maxAge")
	writer.MaxBackups = section.GetInt("maxBackups")
	writer.MaxSize = section.GetInt("maxSize")

	s := strings.TrimSpace(section.GetString("filename"))
	// path.Dir() below only recognize '/' as separator.
	s = strings.ReplaceAll(s, "\\", "/")

	dir := path.Dir(s)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	fn := path.Base(s)
	writer.Filename = path.Join(dir, fn)

	return writer, nil
}

// loadAppenderEncoder loads zapcore.Encoder defined in config file.
// Default value JSON encoder will be used when error occurs.
func loadAppenderEncoder(appender *appenderConfig, appenderSection *viper.Viper) {
	s := strings.ToLower(strings.TrimSpace(appenderSection.GetString("encoderType")))
	var encoder zapcore.Encoder

	if s == "console" {
		encoder = zapcore.NewConsoleEncoder(*appender.encoderConfig)
	} else {
		// all other values are treated as JSON.
		encoder = zapcore.NewJSONEncoder(*appender.encoderConfig)
	}

	appender.encoder = &encoder
}

// loadAppenderLogLevel loads log level defined in config file.
// Default value InfoLevel will be used when error occurs.
func loadAppenderLogLevel(appender *appenderConfig, appenderSection *viper.Viper) {
	atomicLevel := zap.NewAtomicLevel()

	if atomicLevel.UnmarshalText(getLowerBytes(appenderSection, "logLevel")) != nil {
		// undefined value treated as InfoLevel.
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	appender.logLevel = atomicLevel
}

// loadAppenderEncoderConfig returns zapcore.EncoderConfig defined in config file.
// It returns error when encoderConfig entry was missing.
func loadAppenderEncoderConfig(config *viper.Viper, appender *appenderConfig, appenderSection *viper.Viper, appenderName string) error {
	// get value of appender encoderConfig first.
	sectionName, err := getRequiredString(appenderSection, appenderName, "encoderConfig")
	if err != nil {
		return err
	}

	section := config.Sub(sectionName)
	if section == nil {
		return fmt.Errorf("the value of [%s.encoderConfig] is [%s], but the entry was missing",
			appenderName, sectionName)
	}

	encoderConfig := new(zapcore.EncoderConfig)

	// deal with string properties.
	if section.IsSet("callerKey") {
		encoderConfig.CallerKey = section.GetString("callerKey")
	}
	if section.IsSet("functionKey") {
		encoderConfig.FunctionKey = section.GetString("functionKey")
	}
	if section.IsSet("levelKey") {
		encoderConfig.LevelKey = section.GetString("levelKey")
	}
	if section.IsSet("messageKey") {
		encoderConfig.MessageKey = section.GetString("messageKey")
	}
	if section.IsSet("nameKey") {
		encoderConfig.NameKey = section.GetString("nameKey")
	}
	if section.IsSet("stacktraceKey") {
		encoderConfig.StacktraceKey = section.GetString("stacktraceKey")
	}
	if section.IsSet("timeKey") {
		encoderConfig.TimeKey = section.GetString("timeKey")
	}
	if section.IsSet("consoleSeparator") {
		encoderConfig.ConsoleSeparator = section.GetString("consoleSeparator")
	}
	if section.IsSet("lineEnding") {
		encoderConfig.LineEnding = section.GetString("lineEnding")
	}

	// deal with typed properties. must use lower case.
	// no err will be returned so ignore it.
	_ = encoderConfig.EncodeCaller.UnmarshalText(getLowerBytes(section, "callerEncoder"))
	_ = encoderConfig.EncodeDuration.UnmarshalText(getLowerBytes(section, "durationEncoder"))
	_ = encoderConfig.EncodeLevel.UnmarshalText(getLowerBytes(section, "levelEncoder"))
	_ = encoderConfig.EncodeName.UnmarshalText(getLowerBytes(section, "nameEncoder"))

	s := section.GetString("timeEncoder")
	if strings.Index(s, "%") == 0 { // customized format starts from '%'
		s = s[1:]
		// the format string is like "2006-01-02 15:04:05.999999999 -0700 MST",
		// use others will get unpredictable value.
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(s))
		}
	} else {
		_ = encoderConfig.EncodeTime.UnmarshalText(getLowerBytes(section, "timeEncoder"))
	}

	appender.encoderConfig = encoderConfig

	return nil
}

// getLowerBytes returns a byte array from config.
// It gets the string first, then trim it, finally covert it to byte array.
func getLowerBytes(section *viper.Viper, key string) []byte {
	s := strings.ToLower(strings.TrimSpace(section.GetString(key)))
	return []byte(s)
}

// getRequiredString returns a non empty string from config.
// It returns error if the entry doesn't exist, or the value is empty.
func getRequiredString(section *viper.Viper, sectionName string, key string) (string, error) {
	// InConfig() is case sensitive, must be lower case. IsSet() is not case sensitive.
	if section.IsSet(key) {
		value := strings.TrimSpace(section.GetString(key))
		if len(value) > 0 {
			return value, nil
		}

		return "", fmt.Errorf("the value of [%s.%s] is empty", sectionName, key)
	}

	return "", fmt.Errorf("the key [%s.%s] does not exist", sectionName, key)
}
