package healthcheck

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	healthcheckService Service
}

func NewHandler() *Handler {
	return &Handler{healthcheckService: NewService()}
}

func (h *Handler) Check(c echo.Context) error {
	err := h.healthcheckService.Check()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")
}
