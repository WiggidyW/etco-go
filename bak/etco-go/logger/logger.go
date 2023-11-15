package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func MaybeFatal(err error) {
	if err != nil {
		Fatal(err.Error())
	}
}

func Err(msg string) {
	logger.Error(msg)
}

func MaybeErr(err error) {
	if err != nil {
		Err(err.Error())
	}
}

func Warn(msg string) {
	logger.Warn(msg)
}

func MaybeWarn(err error) {
	if err != nil {
		Warn(err.Error())
	}
}

func Info(msg string) {
	logger.Info(msg)
}

func Debug(msg string) {
	logger.Debug(msg)
}
