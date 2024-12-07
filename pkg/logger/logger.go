package logger

import (
	"kaf-interface/internal/orders/config"
	"log/slog"
	"os"

	"github.com/sytallax/prettylog"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func Load(cfg *config.Config) *slog.Logger {

	var logger *slog.Logger

	// логгер в зависимости от environment конфига
	// local - красивый вывод в консоль
	// prod - json-формат в stdout
	switch cfg.Env {
	case envLocal:
		logger = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   true,
			ReplaceAttr: nil,
		}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return logger
}
