package logger

import (
	"go.uber.org/zap"
	"todo_list_go/internal/config"
)

const (
	prodLogEnv = "prod"
)

func Init(loggerCfg config.LoggerConfig) error {
	var cfg zap.Config

	if loggerCfg.LoggerEnv == prodLogEnv {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.DisableStacktrace = true

	baseLogger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger := baseLogger.WithOptions(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)

	return err
}

func Debug(msg string) {
	zap.S().Debug(msg)
}

func Debugf(msg string, args ...interface{}) {
	zap.S().Debugf(msg, args)
}

func Info(msg string) {
	zap.S().Info(msg)
}

func Infof(msg string, args ...interface{}) {
	zap.S().Infof(msg, args)
}

func Warn(msg string) {
	zap.S().Warn(msg)
}

func Warnf(msg string, args ...interface{}) {
	zap.S().Warnf(msg, args)
}

func Error(msg string) {
	zap.S().Error(msg)
}

func Errorf(msg string, args ...interface{}) {
	zap.S().Errorf(msg, args)
}
