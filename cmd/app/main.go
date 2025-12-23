package main

import (
	"github.com/cruffinoni/rimworld-editor/internal/application"
	"github.com/cruffinoni/rimworld-editor/internal/config"
	"github.com/cruffinoni/rimworld-editor/pkg/logging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fallback := &config.Config{
			Logging: config.LoggingConfig{
				Level:  "info",
				Format: "text",
				Output: "stderr",
			},
		}
		logger := logging.NewLogger("app", &fallback.Logging)
		logger.WithError(err).Fatal("Failed to load config")
	}
	logger := logging.NewLogger("app", &cfg.Logging)
	app := application.CreateApplication(logger)
	if err := app.Run(); err != nil {
		logger.WithError(err).Fatal("Application failed")
	}
}
