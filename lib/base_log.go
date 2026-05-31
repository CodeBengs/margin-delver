package lib

import "go.uber.org/zap"

type BaseLog struct {
	log *zap.Logger
}

func NewBaseLog(cfg *AppConfig) *BaseLog {
	var logger *zap.Logger

	if cfg.AppEnv == "local" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	return &BaseLog{
		log: logger,
	}
}

func (b *BaseLog) Logger() *zap.Logger {
	return b.log
}

func (b *BaseLog) SugarLog() *zap.SugaredLogger {
	return b.log.Sugar()
}
