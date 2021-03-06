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
	const SectionName = "appenders"

	if !config.IsSet(SectionName) {
		return nil, nil, fmt.Errorf("missing section [%s]", SectionName)
	}

	appenderNames := config.Get(SectionName).([]interface{})

	// at least one appender is required.
	if len(appenderNames) == 0 {
		return nil, nil, fmt.Errorf("no appender is defined in section [%s]", SectionName)
	}

	appenders := make(map[string]*appenderConfig)
	errorAppenders := make(map[string]error)

	// load appenders from config file and put them into a map to filter duplication.
	for _, v := range appenderNames {
		w := strings.TrimSpace(v.(string))
		appenders[w] = nil
	}

	// load each appender info from corresponding section.
	// the appender name is the section name.
	// put it into error appender list if error happened.
	for appenderName := range appenders {
		if appender, err := loadAppender(config, appenderName); err == nil {
			appenders[appenderName] = appender
		} else {
			errorAppenders[appenderName] = err
		}
	}

	// remove error appenders.
	if len(errorAppenders) > 0 {
		for k := range errorAppenders {
			delete(appenders, k)
		}
	}

	// there should be at least one successful loaded appender.
	if len(appenders) == 0 {
		return nil, errorAppenders, fmt.Errorf("fail to load all %d appenders", len(errorAppenders))
	}

	return appenders, errorAppenders, nil
}

// loadAppender loads each appender according to its name.
// it returns appenderConfig object and error object.
func loadAppender(config *viper.Viper, appenderName string) (*appenderConfig, error) {
	appenderSection := config.Sub(appenderName)
	if appenderSection == nil {
		return nil, fmt.Errorf("cannot find the entry for appender [%s]", appenderName)
	}

	appender := new(appenderConfig)
	appender.name = appenderName

	if err := loadAppenderWriteSyncer(config, appenderSection, appender); err != nil {
		return nil, err
	}
	if err := loadAppenderEncoderConfig(config, appenderSection, appender); err != nil {
		return nil, err
	}

	loadAppenderLogLevel(appender, appenderSection)
	// must be called after loadAppenderEncoderConfig() because Encoder needs EncoderConfig.
	loadAppenderEncoder(appender, appenderSection)

	return appender, nil
}

// loadAppenderWriteSyncer loads WriteSyncer from corresponding appender section.
// it returns error when the entry was missing.
func loadAppenderWriteSyncer(config *viper.Viper, appenderSection *viper.Viper, appender *appenderConfig) error {
	// 'target' is the fixed and required key. the value should not be empty.
	s, err := getRequiredString(appenderSection, appender.name, "target")
	if err != nil {
		return err
	}

	var syncer zapcore.WriteSyncer

	if strings.EqualFold(s, "stdout") {
		syncer = zapcore.AddSync(os.Stdout)
	} else if strings.EqualFold(s, "stderr") {
		syncer = zapcore.AddSync(os.Stderr)
	} else {
		// find lumberjack section.
		section := config.Sub(s)
		if section == nil {
			return fmt.Errorf("the value of [%s.target] is [%s], but the entry was missing", appender.name, s)
		}

		if writer, err := loadLumberjack(section); err != nil {
			return err
		} else {
			syncer = zapcore.AddSync(writer)
		}
	}

	appender.writeSyncer = &syncer

	return nil
}

// loadLumberjack loads lumberjack.Logger as io.Writer from config file.
// it returns error when it failed to create log file path.
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

	filename := path.Base(s)
	writer.Filename = path.Join(dir, filename)

	return writer, nil
}

// loadAppenderEncoderConfig returns zapcore.EncoderConfig defined in config file.
// It returns error when encoderConfig entry was missing.
func loadAppenderEncoderConfig(config *viper.Viper, appenderSection *viper.Viper, appender *appenderConfig) error {
	// get value of appender encoderConfig section name first.
	// 'encoderConfig' is the fixed and required name defined in appender section.
	// its value can be any meaningful word, not required as 'encoderConfig'.
	sectionName, err := getRequiredString(appenderSection, appender.name, "encoderConfig")
	if err != nil {
		return err
	}

	section := config.Sub(sectionName)
	if section == nil {
		return fmt.Errorf("the value of [%s.encoderConfig] is [%s], but the entry was missing",
			appender.name, sectionName)
	}

	encoderConfig := new(zapcore.EncoderConfig)

	// deal with string properties. all keys are optional.
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
	if section.IsSet("encodeDuration") {
		_ = encoderConfig.EncodeDuration.UnmarshalText(getLowerBytes(section, "encodeDuration"))
	}
	if section.IsSet("encodeLevel") {
		_ = encoderConfig.EncodeLevel.UnmarshalText(getLowerBytes(section, "encodeLevel"))
	}
	if section.IsSet("encodeName") {
		_ = encoderConfig.EncodeName.UnmarshalText(getLowerBytes(section, "encodeName"))
	}
	if section.IsSet("encodeTime") {
		s := section.GetString("encodeTime")
		if strings.Index(s, "%") == 0 { // customized format starts from '%'
			s = s[1:]
			// the format string is like "2006-01-02 15:04:05.999999999 -0700 MST",
			// use others will get unpredictable value.
			encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(s))
			}
		} else {
			_ = encoderConfig.EncodeTime.UnmarshalText(getLowerBytes(section, "encodeTime"))
		}
	}

	// no idea why it throws error when we don't set this filed.
	_ = encoderConfig.EncodeCaller.UnmarshalText(getLowerBytes(section, "encodeCaller"))

	appender.encoderConfig = encoderConfig

	return nil
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

// loadAppenderEncoder loads zapcore.Encoder defined in config file.
// default value JSON encoder will be used when error occurs.
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
