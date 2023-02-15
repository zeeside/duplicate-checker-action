package main

import (
	log "github.com/sirupsen/logrus"
)

type ActionLogger struct {
	logger *log.Logger
}

func NewActionLogger(cfg *Config) *ActionLogger {
	loggerInstance := log.New()
	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = log.InfoLevel
	}

	loggerInstance.SetLevel(level)

	return &ActionLogger{
		logger: loggerInstance,
	}
}
