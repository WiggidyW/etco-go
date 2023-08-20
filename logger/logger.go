package logger

import (
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

func Warn(err error) {
	if err != nil {
		Logger.Warn(err.Error())
	}
}

func Err(err error) {
	if err != nil {
		Logger.Error(err.Error())
	}
}

func Fatal(err error) {
	if err != nil {
		Logger.Fatal(err.Error())
	}
}
