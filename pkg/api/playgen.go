package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/meyskens/recent-beater/pkg/bplist"

	"github.com/labstack/echo/v4"
	"github.com/meyskens/recent-beater/pkg/scoresaber"
)

var alphaNum = regexp.MustCompile("[^a-zA-Z0-9]+")

type Mode string

const (
	ModeRecent Mode = "recent"
	ModeTop    Mode = "top"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/playlist/", h.GeneratePlaylist)
		e.GET("/playlist/:id/:pages/:mode/:filename", h.GeneratePlaylistWithID)
	})
}

func (h *HTTPHandler) GeneratePlaylistWithID(c echo.Context) error {
	pid := c.Param("id")

	if pid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "player id not specified"})
	}

	mode := ModeRecent
	if c.Param("mode") == "top" {
		mode = ModeTop
	}

	amount := 1
	if i, err := strconv.ParseInt(c.Param("pages"), 10, 64); err == nil && i > 0 {
		amount = int(i)
	}

	return h.handlePlaygen(c, pid, amount, mode)
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
	if i, err := strconv.ParseInt(c.QueryParam("pages"), 10, 64); err == nil && i > 0 {
		amount = int(i)
	}

	mode := ModeRecent
	if c.QueryParam("mode") == "top" {
		mode = ModeTop
	}

	if c.QueryParam("oneclick") == "true" {
		profile, err := scoresaber.GetProfile(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		// bsplaylist has no support for GET parameters
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("bsplaylist://playlist/https://recentbeat.com/playlist/%s/%d/%s/BEAT_%s_%s_%d.bplist\n", id, amount, mode, alphaNum.ReplaceAllString(profile.PlayerInfo.PlayerName, ""), mode, time.Now().Unix()))
	}

	return h.handlePlaygen(c, id, amount, mode)
}

func (h *HTTPHandler) handlePlaygen(c echo.Context, id string, amount int, mode Mode) error {

	if amount > 10 {
		amount = 10 // sorry i am not blowing my rate limit
	}

	var scores []scoresaber.Score
	for i := 0; i < amount; i++ {
		var s []scoresaber.Score
		var err error
		if mode == ModeRecent {
			s, err = scoresaber.GetRecentScores(id, i+1)

		} else if mode == ModeTop {
			s, err = scoresaber.GetTopScores(id, i+1)
		}
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
	pl.PlaylistTitle = fmt.Sprintf("BEAT %s %s", mode, profile.PlayerInfo.PlayerName)
	pl.PlaylistAuthor = "Recent Beater"

	for _, score := range scores {
		pl.AddSong(score.ToBPlistSong())
	}

	c.Response().Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.bplist\"", pl.PlaylistTitle))

	return c.JSON(http.StatusOK, pl)
}
