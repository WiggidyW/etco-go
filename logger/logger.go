package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger   *zap.Logger   = nil
	LoggerMu *sync.RWMutex = &sync.RWMutex{}
)

func InitLoggerCrashOnError() {
	LoggerMu.Lock()
	defer LoggerMu.Unlock()

	var err error
	// Logger, err = zap.NewProduction()
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func loggerRLockIfNil() {
	if logger == nil {
		LoggerMu.RLock()
		defer LoggerMu.RUnlock()
	}
}

func Warn(err error) {
	if err != nil {
		loggerRLockIfNil()
		logger.Warn(err.Error())
	}
}

func Err(err error) {
	if err != nil {
		loggerRLockIfNil()
		logger.Error(err.Error())
	}
}

func Fatal(err error) {
	if err != nil {
		loggerRLockIfNil()
		logger.Fatal(err.Error())
	}
}

func Debug(msg string) {
	loggerRLockIfNil()
	logger.Debug(msg)
}

func Info(msg string) {
	loggerRLockIfNil()
	logger.Info(msg)
}
