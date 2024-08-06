package server

import (
	"school/internal/api"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
)

func RegisteringSwaggerRoutes(s Server) error {
	oapi3, err := api.GetSwagger()
	if err != nil {
		return err
	}

	spec, err := oapi3.MarshalJSON()
	if err != nil {
		return err
	}

	var swaggerInfo = &swag.Spec{
		Version:          "",
		Host:             "",
		BasePath:         "",
		Schemes:          []string{},
		Title:            "",
		Description:      "",
		InfoInstanceName: "",
		SwaggerTemplate:  string(spec),
		LeftDelim:        "",
		RightDelim:       "",
	}
	swag.Register("swagger", swaggerInfo)

	s.engine.GET("/swagger/*any", echoSwagger.WrapHandler)
	return nil
}
