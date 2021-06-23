package cfzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// the logger that should be created according to configuration.
	logger *zap.Logger

	// default logger in case failed to create logger from config file.
	defaultLogger *zap.Logger
)

// GetLogger returns a logger according to the config file.
// If 'createNew' is true, then trying to return the exist logger created before.
// you can set config file information by calling SetLoggerConfig() before call this function.
// In theory, even with an error, the returned logger will not be nil.
func GetLogger(createNew bool) (*zap.Logger, error) {
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

	if !createNew && logger != nil {
		return logger, nil
	}

	//config, err := readConfigFile(lastConfigFileWithoutExt, lastConfigFileExit, lastConfigPaths...)
	config, err := readConfigFile(nil)
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
		cores[i] = zapcore.NewCore(*appender.Encoder, *appender.WriteSyncer, appender.LogLevel)
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

	return logger, nil
}
