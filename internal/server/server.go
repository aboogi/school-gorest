package server

import (
	"context"
	"fmt"
	"school/internal/config"
	"school/internal/storage"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	host    string
	port    int
	engine  *echo.Echo
	logger  *zap.Logger
	storage *storage.Connection
}

func New(cfg config.Config, logger *zap.Logger, storage *storage.Connection) Server {
	engine := echo.New()
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	engine.Use(middleware.CORS())
	engine.Use(middleware.Gzip())
	engine.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "custom timeout error message returns to client",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Info(c.Path())
		},
		Timeout: 60 * time.Second,
	}))

	s := Server{
		host:    cfg.Host,
		port:    cfg.Port,
		engine:  engine,
		logger:  logger,
		storage: storage,
	}

	RegisteringMetricsRoute(s)
	RegisteringSwaggerRoutes(s)
	RegisteringAPIRoutes(s, "api/v1")

	rr := s.engine.Router().Routes()
	for _, v := range rr {
		fmt.Println(v.Method, v.Path)
	}

	return s
}

func (s Server) Start(ctx context.Context) {
	go func() {
		if err := s.engine.Start(fmt.Sprintf(":%d", s.port)); err != nil {
			s.logger.Sugar().Fatalf("shutting down the server: %w", err)
		}
	}()
	<-ctx.Done()

	err := s.engine.Shutdown(context.Background())
	if err != nil {
		s.logger.Sugar().Fatalf("error on shutdown apiserver: %w", err)
	}
}
