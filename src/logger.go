package cfzap

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// the logger that should be created according to configuration.
	logger *zap.Logger

	// default logger in case failed to create logger from config file.
	defaultLogger *zap.Logger

	// the ConfigOption was used last time.
	lastConfigOption *ConfigOption

	lock sync.Mutex
)

// GetLogger returns a logger according to the config file.
// If 'createNew' is true, then trying to return the exist logger created before.
// In theory, even with an error, the returned logger will not be nil.
func GetLogger(configOption *ConfigOption) (*zap.Logger, error) {
	lock.Lock()
	defer lock.Unlock()

	// create default logger when there's no configured logger to use.
	if defaultLogger == nil {
		var e error
		if defaultLogger, e = zap.NewDevelopment(zap.AddCaller()); e != nil {
			// such simple code should never go wrong. if it really happens, we have to quit.
			panic(e)
		}
	}

	defer func() {
		if defaultLogger != nil {
			_ = defaultLogger.Sync()
		}
		if logger != nil {
			_ = logger.Sync()
		}
	}()

	if configOption == nil { // using default value if it is not provided.
		configOption = NewConfigOption()
	}

	if !shouldCreateNewLogger(configOption) {
		return logger, nil
	}

	config, err := readConfigFile(configOption)
	if err != nil {
		defaultLogger.Warn("fail to load logger config: " + err.Error())
		return defaultLogger, err
	}

	appenders, errors, err := loadAppenders(config)

	if err != nil {
		return defaultLogger, err
	}

	cores := make([]zapcore.Core, len(appenders))
	i := 0
	for _, appender := range appenders {
		cores[i] = zapcore.NewCore(*appender.encoder, *appender.writeSyncer, appender.logLevel)
		i++
	}

	for k, v := range errors {
		defaultLogger.Warn("fail to load appender [" + k + "]: " + v.Error())
	}

	if logger != nil { // flush old logger before creating a new one.
		_ = logger.Sync()
	}

	options := loadLogOptions(config)
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, options...)

	lastConfigOption = configOption

	return logger, nil
}

// shouldCreateNewLogger check to see if we should create a new logger object.
func shouldCreateNewLogger(configOption *ConfigOption) bool {
	if logger == nil { // we don't have a created logger yet.
		return true
	}
	if configOption.CreateNew { // the option does require create a new one.
		return true
	}

	// we have to create a new one when the new option is diffrent compare to the last one.
	return !lastConfigOption.equal(configOption)
}
