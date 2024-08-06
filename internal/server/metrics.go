package server

import "github.com/labstack/echo-contrib/echoprometheus"

func RegisteringMetricsRoute(s Server) {
	s.engine.GET("/metrics", echoprometheus.NewHandler())
}
