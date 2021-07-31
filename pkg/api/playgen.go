package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/meyskens/recent-beater/pkg/bplist"

	"github.com/labstack/echo/v4"
	"github.com/meyskens/recent-beater/pkg/scoresaber"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/playlist/", h.GeneratePlaylist)
		e.GET("/playlist/:id/:pages/:filename", h.GeneratePlaylistWithID)
	})
}

func (h *HTTPHandler) GeneratePlaylistWithID(c echo.Context) error {
	pid := c.Param("id")

	if pid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "player id not specified"})
	}

	amount := 1
	if i, err := strconv.ParseInt(c.Param("pages"), 10, 64); err != nil && i > 0 {
		amount = int(i)
	}

	return h.handlePlaygen(c, pid, amount)
}

func (h *HTTPHandler) GeneratePlaylist(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "URL not specified"})
	}

	id, err := scoresaber.ExtractPlayerIDFromURL(url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	amount := 1
	if i, err := strconv.ParseInt(c.QueryParam("pages"), 10, 64); err != nil && i > 0 {
		amount = int(i)
	}

	if c.QueryParam("oneclick") == "true" {
		profile, err := scoresaber.GetProfile(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		// bsplaylist has no support for GET parameters
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("bsplaylist://playlist/https://recentbeat.com/playlist/%s/%d/BEAT_%s_%s.bplist\n", id, amount, profile.PlayerInfo.PlayerName, time.Now().UTC().String()))
	}

	return h.handlePlaygen(c, id, amount)
}

func (h *HTTPHandler) handlePlaygen(c echo.Context, id string, amount int) error {
	if id == "76561198407185197" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No Dirk, thy shall not use my own tools against me"})
	}

	var scores []scoresaber.Score
	for i := 0; i < amount; i++ {
		s, err := scoresaber.GetRecentScores(id, i+1)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		if len(s) == 0 {
			break
		}

		scores = append(scores, s...)
	}

	profile, err := scoresaber.GetProfile(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	pl := bplist.NewPlaylist()
	pl.PlaylistTitle = fmt.Sprintf("BEAT %s", profile.PlayerInfo.PlayerName)
	pl.PlaylistAuthor = "Recent Beater"

	for _, score := range scores {
		pl.AddSong(score.ToBPlistSong())
	}

	c.Response().Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.bplist\"", pl.PlaylistTitle))

	return c.JSON(http.StatusOK, pl)
}
