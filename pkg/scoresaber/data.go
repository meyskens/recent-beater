package scoresaber

import "time"

type scoresAPIData struct {
	Scores []Score `json:"scores"`
}

// Score is a single score a player made on scoresaber
type Score struct {
	Rank              int       `json:"rank"`
	ScoreID           int       `json:"scoreId"`
	Score             int       `json:"score"`
	UnmodififiedScore int       `json:"unmodififiedScore"`
	Mods              string    `json:"mods"`
	PP                float64   `json:"pp"`
	Weight            float64   `json:"weight"`
	TimeSet           time.Time `json:"timeSet"`
	LeaderboardID     int       `json:"leaderboardId"`
	SongHash          string    `json:"songHash"`
	SongName          string    `json:"songName"`
	SongSubName       string    `json:"songSubName"`
	SongAuthorName    string    `json:"songAuthorName"`
	LevelAuthorName   string    `json:"levelAuthorName"`
	Difficulty        int       `json:"difficulty"`
	DifficultyRaw     string    `json:"difficultyRaw"`
	MaxScore          int       `json:"maxScore"`
}

type Profile struct {
	PlayerInfo struct {
		PlayerID    string        `json:"playerId"`
		PlayerName  string        `json:"playerName"`
		Avatar      string        `json:"avatar"`
		Rank        int           `json:"rank"`
		CountryRank int           `json:"countryRank"`
		Pp          float64       `json:"pp"`
		Country     string        `json:"country"`
		Role        interface{}   `json:"role"`
		Badges      []interface{} `json:"badges"`
		History     string        `json:"history"`
		Permissions int           `json:"permissions"`
		Inactive    int           `json:"inactive"`
		Banned      int           `json:"banned"`
	} `json:"playerInfo"`
	ScoreStats struct {
		TotalScore            int     `json:"totalScore"`
		TotalRankedScore      int     `json:"totalRankedScore"`
		AverageRankedAccuracy float64 `json:"averageRankedAccuracy"`
		TotalPlayCount        int     `json:"totalPlayCount"`
		RankedPlayCount       int     `json:"rankedPlayCount"`
	} `json:"scoreStats"`
}
