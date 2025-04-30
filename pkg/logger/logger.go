package logger

import "go.uber.org/zap"

func Init() error {
	config := zap.NewDevelopmentConfig()

	config.DisableStacktrace = true

	baseLogger, err := config.Build()
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
