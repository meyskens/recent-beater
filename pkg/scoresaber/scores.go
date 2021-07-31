package scoresaber

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const HOST = "https://new.scoresaber.com/api"

// ExtractPlayerIDFromURL gets you the player ID for a given scoresaber URL
func ExtractPlayerIDFromURL(u string) (string, error) {
	uri, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	parts := strings.Split(uri.Path, "/")
	var id string
	for i, part := range parts {
		if part == "u" {
			id = parts[i+1]
		}
	}

	if id == "" {
		return "", errors.New("ID part not found")
	}

	return id, nil
}

// GetRecentScores gets the recent scores of a given player
func GetRecentScores(playerID string, page int) ([]Score, error) {
	resp, err := http.Get(fmt.Sprintf("%s/player/%s/scores/recent/%d", HOST, playerID, page))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data scoresAPIData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Scores, nil
}

// GetTopScores gets the top scores of a given player
func GetTopScores(playerID string, page int) ([]Score, error) {
	resp, err := http.Get(fmt.Sprintf("%s/player/%s/scores/top/%d", HOST, playerID, page))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data scoresAPIData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Scores, nil
}

// GetProfile gets the profile info of a given player
func GetProfile(playerID string) (*Profile, error) {
	resp, err := http.Get(fmt.Sprintf("%s/player/%s/full", HOST, playerID))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Profile
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
