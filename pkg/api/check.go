package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/meyskens/recent-beater/pkg/scoresaber"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/check/", h.CheckURL)
	})
}

func (h *HTTPHandler) CheckURL(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "URL not specified"})
	}

	id, err := scoresaber.ExtractPlayerIDFromURL(url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err = scoresaber.GetProfile(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
