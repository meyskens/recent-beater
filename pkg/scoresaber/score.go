package scoresaber

import (
	"strings"

	"github.com/meyskens/recent-beater/pkg/bplist"
)

func (s *Score) ToBPlistSong() bplist.Song {
	var difficulties []bplist.Difficulty
	difParts := strings.Split(s.DifficultyRaw, "_")
	if len(difParts) > 2 {
		difficulties = append(difficulties, bplist.Difficulty{
			Characteristic: "Standard",
			Name:           difParts[1],
		})
	}
	return bplist.Song{
		Hash:         s.SongHash,
		SongName:     s.SongName,
		Difficulties: difficulties,
	}
}
