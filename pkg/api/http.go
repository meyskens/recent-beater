package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var registers []func(e *echo.Echo, h *HTTPHandler)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "OwO")
		})
	})
}

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	for _, regFn := range registers {
		regFn(e, h)
	}
}
