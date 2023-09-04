package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	Logger   *zap.Logger = nil
	LoggerMu             = &sync.RWMutex{}
)

func InitLoggerCrashOnError() {
	LoggerMu.Lock()
	defer LoggerMu.Unlock()

	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func loggerRLockIfNil() {
	if Logger == nil {
		LoggerMu.RLock()
		defer LoggerMu.RUnlock()
	}
}

func Warn(err error) {
	if err != nil {
		loggerRLockIfNil()
		Logger.Warn(err.Error())
	}
}

func Err(err error) {
	if err != nil {
		loggerRLockIfNil()
		Logger.Error(err.Error())
	}
}

func Fatal(err error) {
	if err != nil {
		loggerRLockIfNil()
		Logger.Fatal(err.Error())
	}
}
