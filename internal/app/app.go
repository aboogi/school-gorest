package app

import (
	"context"
	"schoolmat/internal/config"
	"schoolmat/internal/server"
	"schoolmat/internal/storage"

	"go.uber.org/zap"
)

type App struct {
	config config.Config
	logger *zap.Logger
}

func New(cfg config.Config, logger *zap.Logger) App {
	return App{
		config: cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) {
	connection, err := storage.NewConnection(ctx, a.config.DB)
	if err != nil {
		a.logger.Sugar().Error("error connecting to database: %w", err)
		return
	}

	server.New(a.config, a.logger, &connection).Start(ctx)
}
