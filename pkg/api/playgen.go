package api

import (
	"encoding/base64"
	"encoding/json"
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

type options struct {
	URL        string `json:"url"`
	ID         string `json:"id"`
	Mode       mode   `json:"mode"`
	Pages      int    `json:"pages"`
	OneClick   bool   `json:"oneClick"`
	FilterNoPP bool   `json:"filterNoPP"`
}

func newOptions() options {
	return options{
		URL:        "",
		ID:         "",
		Mode:       modeRecent,
		Pages:      1,
		OneClick:   false,
		FilterNoPP: false,
	}
}

type mode string

const (
	modeRecent mode = "recent"
	modeTop    mode = "top"
)

func init() {
	registers = append(registers, func(e *echo.Echo, h *HTTPHandler) {
		e.GET("/playlist/", h.GeneratePlaylist)
		e.GET("/playlist/:data/:filename", h.GeneratePlaylistWithID)
	})
}

func (h *HTTPHandler) GeneratePlaylistWithID(c echo.Context) error {
	data, err := base64.StdEncoding.DecodeString(c.Param("data"))
	if err != nil || len(data) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data payload"})
	}

	opts := newOptions()
	err = json.Unmarshal(data, &opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data payload"})
	}

	return h.handlePlaygen(c, opts)
}

func (h *HTTPHandler) GeneratePlaylist(c echo.Context) error {
	opts := newOptions()

	opts.URL = c.QueryParam("url")
	if opts.URL == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "URL not specified"})
	}

	id, err := scoresaber.ExtractPlayerIDFromURL(opts.URL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	opts.ID = id

	if i, err := strconv.ParseInt(c.QueryParam("pages"), 10, 64); err == nil && i > 0 {
		opts.Pages = int(i)
	}

	if c.QueryParam("filterNoPP") == "true" {
		opts.FilterNoPP = true
	}

	if c.QueryParam("mode") == "top" {
		opts.Mode = modeTop
	}

	if c.QueryParam("oneclick") == "true" {
		opts.OneClick = true

		profile, err := scoresaber.GetProfile(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		data, err := json.Marshal(opts)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		// bsplaylist has no support for GET parameters
		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("bsplaylist://playlist/https://recentbeat.com/playlist/%s/BEAT_%s_%s_%d.bplist\n", base64.StdEncoding.EncodeToString(data), alphaNum.ReplaceAllString(profile.PlayerInfo.PlayerName, ""), opts.Mode, time.Now().Unix()))
	}

	return h.handlePlaygen(c, opts)
}

func (h *HTTPHandler) handlePlaygen(c echo.Context, opts options) error {
	if opts.ID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID not specified"})
	}

	if opts.Pages > 10 {
		opts.Pages = 10 // sorry i am not blowing my rate limit
	}

	var scores []scoresaber.Score
	for i := 0; i < opts.Pages; i++ {
		var s []scoresaber.Score
		var err error
		if opts.Mode == modeRecent {
			s, err = scoresaber.GetRecentScores(opts.ID, i+1)

		} else if opts.Mode == modeTop {
			s, err = scoresaber.GetTopScores(opts.ID, i+1)
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		if len(s) == 0 {
			break
		}

		scores = append(scores, s...)
	}

	if opts.FilterNoPP {
		var newScores []scoresaber.Score
		for _, score := range scores {
			if score.PP > 0 {
				newScores = append(newScores, score)
			}
		}

		scores = newScores
	}

	profile, err := scoresaber.GetProfile(opts.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	pl := bplist.NewPlaylist()
	pl.PlaylistTitle = fmt.Sprintf("BEAT %s %s", opts.Mode, profile.PlayerInfo.PlayerName)
	pl.PlaylistAuthor = "Recent Beater"

	for _, score := range scores {
		pl.AddSong(score.ToBPlistSong())
	}

	c.Response().Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.bplist\"", pl.PlaylistTitle))

	return c.JSON(http.StatusOK, pl)
}
