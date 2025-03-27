package healthcheck

import "github.com/labstack/echo/v4"

func NewGroup(g *echo.Group) {
	handler := NewHandler()

	group := g.Group("/healthcheck")
	group.GET("", handler.Check)
}
