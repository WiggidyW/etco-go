package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger   *zap.Logger = nil
	loggerMu *sync.Mutex = &sync.Mutex{}
)

func getLogger() *zap.Logger {
	if logger == nil {
		loggerMu.Lock()
		defer loggerMu.Unlock()

		if logger == nil {
			var err error
			logger, err = zap.NewDevelopment()
			if err != nil {
				panic(err)
			}
		}
	}
	return logger
}

func InitLoggerCrashOnError() {
	getLogger()
}

func Warn(err error) {
	if err != nil {
		getLogger().Warn(err.Error())
	}
}

func Err(err error) {
	if err != nil {
		getLogger().Error(err.Error())
	}
}

func Fatal(err error) {
	if err != nil {
		getLogger().Fatal(err.Error())
	}
}

func Debug(msg string) {
	getLogger().Debug(msg)
}

func Info(msg string) {
	getLogger().Info(msg)
}
