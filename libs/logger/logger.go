package logger

import "go.uber.org/zap"

var rootLogger *zap.Logger

func Init(isDebug bool) {
	rootLogger = New(isDebug)
}

func New(isDebug bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if isDebug {
		logger, err = zap.NewDevelopment()
	} else {
		config := zap.NewProductionConfig()
		config.DisableCaller = true
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = config.Build()
	}
	if err != nil {
		panic(err)
	}
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	rootLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	rootLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	rootLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	rootLogger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	rootLogger.Fatal(msg, fields...)
}
