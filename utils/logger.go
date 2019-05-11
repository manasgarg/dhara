package utils

import (
	"go.uber.org/zap"
)

var (
	Logger  *zap.Logger
	SLogger *zap.SugaredLogger
)

func InitLogger() {
	Logger, _ = zap.NewDevelopment()
	SLogger = Logger.Sugar()
}
