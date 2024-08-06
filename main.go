package main

import (
	"context"
	"os"
	"os/signal"
	"schoolmat/internal/app"
	"schoolmat/internal/config"
	"schoolmat/internal/logger"
	"schoolmat/sql"

	"github.com/caarlos0/env/v7"
)

func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New()
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err = sql.SetMigration(ctx, "postgres", cfg.DB)
	if err != nil {
		panic(err)
	}

	app := app.New(cfg, logger)
	app.Run(ctx)
}
