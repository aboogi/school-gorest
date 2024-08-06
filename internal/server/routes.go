package server

import "school/internal/api"

func RegisteringAPIRoutes(s Server, baseURL string) {
	r := Routes{
		storage: s.storage,
		logger:  s.logger,
	}

	api.RegisterHandlersWithBaseURL(s.engine, &r, baseURL)
}
