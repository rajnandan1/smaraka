package logger

import (
	"github.com/rajnandan1/smaraka/constants"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func StartLogger(env string) {
	var logger *zap.Logger
	if env == constants.EnvProduction {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync() // flushes buffer, if any
	sugar = logger.Sugar()
	sugar.Info("Logger started")
}

// log info
func LogInfo(message string, fields ...interface{}) {
	if fields == nil {
		sugar.Info(message)
		return
	}
	sugar.Info(message, fields)
}

// log error
func LogError(message string, fields ...interface{}) {
	if fields == nil {
		sugar.Error(message)
		return
	}
	sugar.Errorln(message, fields)
}
