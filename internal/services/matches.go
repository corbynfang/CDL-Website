package services

import (
	"context"
	"time"

	"github.com/corbynfang/CDL-Website/internal/store"
)

// MatchDetail is the full match response including per-map player scoreboards.
type MatchDetail struct {
	Match MatchInfo `json:"match"`
	Maps  []MapInfo `json:"maps"`
}

type MatchInfo struct {
	ID             uint      `json:"id"`
	TournamentID   uint      `json:"tournament_id"`
	TournamentName string    `json:"tournament_name"`
	TournamentSlug string    `json:"tournament_slug"`
	SeasonName     string    `json:"season_name"`
	GameCode       string    `json:"game_code"`
	Team1ID        uint      `json:"team1_id"`
	Team1Name      string    `json:"team1_name"`
	Team1Abbr      string    `json:"team1_abbr"`
	Team1Logo      string    `json:"team1_logo"`
	Team2ID        uint      `json:"team2_id"`
	Team2Name      string    `json:"team2_name"`
	Team2Abbr      string    `json:"team2_abbr"`
	Team2Logo      string    `json:"team2_logo"`
	Team1Score     int       `json:"team1_score"`
	Team2Score     int       `json:"team2_score"`
	WinnerID       *uint     `json:"winner_id"`
	MatchDate      time.Time `json:"match_date"`
	Format         string    `json:"format"`
	BracketRound   string    `json:"bracket_round"`
}

type MapInfo struct {
	MapNumber   int          `json:"map_number"`
	MapName     string       `json:"map_name"`
	Mode        string       `json:"mode"`
	Score1      int          `json:"score_1"`
	Score2      int          `json:"score_2"`
	WinnerID    *uint        `json:"winner_id"`
	DurationSec int          `json:"duration_sec"`
	Played      bool         `json:"played"`
	Team1Stats  []PlayerStat `json:"team1_stats"`
	Team2Stats  []PlayerStat `json:"team2_stats"`
}

type PlayerStat struct {
	PlayerID        uint    `json:"player_id"`
	Gamertag        string  `json:"gamertag"`
	Kills           int     `json:"kills"`
	Deaths          int     `json:"deaths"`
	KDRatio         float64 `json:"kd_ratio"`
	Damage          int     `json:"damage"`
	Assists         int     `json:"assists"`
	BPRating        float64 `json:"bp_rating"`
	HillTime        int     `json:"hill_time"`
	SndRounds       int     `json:"snd_rounds"`
	PlantCount      int     `json:"plant_count"`
	DefuseCount     int     `json:"defuse_count"`
	FirstBloodCount int     `json:"first_blood_count"`
	FirstDeathCount int     `json:"first_death_count"`
	NonTradedKills  int     `json:"non_traded_kills"`
	HighestStreak   int     `json:"highest_streak"`
	DataQualityNote string  `json:"data_quality_note,omitempty"`
}

// TournamentDetail is the enriched tournament response including derived team count and format.
type MatchService struct {
	matches store.MatchStore
}

func NewMatchService(matches store.MatchStore) *MatchService {
	return &MatchService{matches: matches}
}

func (ms *MatchService) GetMatchDetail(ctx context.Context, id int) (*MatchDetail, error) {
	match, err := ms.matches.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	matchMaps, err := ms.matches.GetMaps(ctx, id)
	if err != nil {
		return nil, err
	}

	statRows, err := ms.matches.GetStatRows(ctx, id)
	if err != nil {
		return nil, err
	}

	statsByMap := map[int][]store.MatchStatRow{}
	for _, s := range statRows {
		statsByMap[s.MapNumber] = append(statsByMap[s.MapNumber], s)
	}

	maps := make([]MapInfo, 0, len(matchMaps))
	for _, mm := range matchMaps {
		out := MapInfo{
			MapNumber:   mm.MapNumber,
			MapName:     mm.MapName,
			Mode:        mm.Mode,
			Score1:      mm.Score1,
			Score2:      mm.Score2,
			WinnerID:    mm.WinnerID,
			DurationSec: mm.DurationSec,
			Played:      mm.Played,
			Team1Stats:  []PlayerStat{},
			Team2Stats:  []PlayerStat{},
		}
		for _, s := range statsByMap[mm.MapNumber] {
			p := PlayerStat{
				PlayerID:        s.PlayerID,
				Gamertag:        s.Gamertag,
				Kills:           s.Kills,
				Deaths:          s.Deaths,
				KDRatio:         s.KDRatio,
				Damage:          s.Damage,
				Assists:         s.Assists,
				BPRating:        s.BPRating,
				HillTime:        s.HillTime,
				SndRounds:       s.SndRounds,
				PlantCount:      s.PlantCount,
				DefuseCount:     s.DefuseCount,
				FirstBloodCount: s.FirstBloodCount,
				FirstDeathCount: s.FirstDeathCount,
				NonTradedKills:  s.NonTradedKills,
				HighestStreak:   s.HighestStreak,
				DataQualityNote: s.DataQualityNote,
			}
			if s.TeamID == match.Team1ID {
				out.Team1Stats = append(out.Team1Stats, p)
			} else {
				out.Team2Stats = append(out.Team2Stats, p)
			}
		}
		maps = append(maps, out)
	}

	return &MatchDetail{
		Match: MatchInfo{
			ID:             match.ID,
			TournamentID:   match.TournamentID,
			TournamentName: match.Tournament.Name,
			TournamentSlug: match.Tournament.Slug,
			SeasonName:     match.Tournament.Season.Name,
			GameCode:       match.Tournament.Season.GameCode,
			Team1ID:        match.Team1ID,
			Team1Name:      match.Team1.Name,
			Team1Abbr:      match.Team1.Abbreviation,
			Team1Logo:      match.Team1.LogoURL,
			Team2ID:        match.Team2ID,
			Team2Name:      match.Team2.Name,
			Team2Abbr:      match.Team2.Abbreviation,
			Team2Logo:      match.Team2.LogoURL,
			Team1Score:     match.Team1Score,
			Team2Score:     match.Team2Score,
			WinnerID:       match.WinnerID,
			MatchDate:      match.MatchDate,
			Format:         match.Format,
			BracketRound:   match.BracketRound,
		},
		Maps: maps,
	}, nil
}
