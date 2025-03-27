package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/polyclient/polyclient/api/features/healthcheck"
	"github.com/polyclient/polyclient/gui"
)

type Router struct {
	echo *echo.Echo
}

func NewRouter() *Router {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())

	e.StaticFS("/", gui.DistDirFS)

	api := e.Group("/api/v1")
	healthcheck.NewGroup(api)

	return &Router{echo: e}
}
