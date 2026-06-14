package services

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
	"golang.org/x/sync/errgroup"
)

type playerRes struct {
	player *models.Player
	err    error
}

type tournamentRes struct {
	val []models.PlayerTournamentStats
	err error
}

type splitRes struct {
	val []store.ModeKDSplit
	err error
}

type PlayerService struct {
	store store.PlayerStore
}

type PlayerKDStats struct {
	PlayerID       int                 `json:"player_id"`
	Gamertag       string              `json:"gamertag"`
	AvatarURL      string              `json:"avatar_url"`
	TotalKills     int                 `json:"total_kills"`
	TotalDeaths    int                 `json:"total_deaths"`
	TotalAssists   int                 `json:"total_assists"`
	AvgKD          float64             `json:"avg_kd"`
	HpKDRatio      float64             `json:"hp_kd_ratio"`
	SndKDRatio     float64             `json:"snd_kd_ratio"`
	ControlKDRatio float64             `json:"control_kd_ratio"`
	Tournaments    []TournamentKDEntry `json:"tournament_stats"`
}

type TournamentKDEntry struct {
	TournamentID   uint    `json:"tournament_id"`
	TournamentName string  `json:"tournament_name"`
	Kills          int     `json:"kills"`
	Deaths         int     `json:"deaths"`
	Assists        int     `json:"assists"`
	KDRatio        float64 `json:"kd_ratio"`
	MapsPlayed     int     `json:"maps_played"`
}

type PlayerMatchHistory struct {
	PlayerID int          `json:"player_id"`
	Events   []MatchEvent `json:"events"`
	Total    int          `json:"total"`
}

type MatchEvent struct {
	Event        string        `json:"event"`
	Year         int           `json:"year"`
	TournamentID uint          `json:"tournament_id"`
	Matches      []MatchResult `json:"matches"`
}

type MatchResult struct {
	MatchID      uint    `json:"match_id"`
	Date         string  `json:"date"`
	Opponent     string  `json:"opponent"`
	OpponentAbbr string  `json:"opponent_abbr"`
	Result       string  `json:"result"`
	KD           float64 `json:"kd"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
}

func CalculateKD(kills, deaths int) float64 {
	if deaths == 0 {
		return 0
	}
	return float64(kills) / float64(deaths)
}

func NewPlayerService(s store.PlayerStore) *PlayerService {
	return &PlayerService{store: s}
}

func (ps *PlayerService) List(ctx context.Context, search string, limit, offset int) ([]models.Player, int64, error) {
	return ps.store.List(ctx, search, limit, offset)
}

func (ps *PlayerService) GetByID(ctx context.Context, id int) (*models.Player, error) {
	return ps.store.GetByID(ctx, id)
}

func (ps *PlayerService) GetMatchStats(ctx context.Context, playerID int) ([]models.PlayerMatchStats, error) {
	return ps.store.ListMatchStats(ctx, playerID)
}

func (ps *PlayerService) GetKDStats(ctx context.Context, playerID int) (*PlayerKDStats, error) {
	g, ctx := errgroup.WithContext(ctx)

	var player *models.Player
	var tournamentStats []models.PlayerTournamentStats
	var splits []store.ModeKDSplit

	g.Go(func() error {
		var err error
		player, err = ps.store.GetByID(ctx, playerID)
		return err
	})

	g.Go(func() error {
		var err error
		tournamentStats, err = ps.store.ListTournamentStats(ctx, playerID)
		return err
	})

	g.Go(func() error {
		var err error
		splits, err = ps.store.ListModeKDSplits(ctx, playerID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var totalKills, totalDeaths, totalAssists int
	entries := make([]TournamentKDEntry, 0, len(tournamentStats))
	for _, stat := range tournamentStats {
		totalKills += stat.TotalKills
		totalDeaths += stat.TotalDeaths
		totalAssists += stat.TotalAssists
		entries = append(entries, TournamentKDEntry{
			TournamentID:   stat.TournamentID,
			TournamentName: stat.Tournament.Name,
			Kills:          stat.TotalKills,
			Deaths:         stat.TotalDeaths,
			Assists:        stat.TotalAssists,
			KDRatio:        CalculateKD(stat.TotalKills, stat.TotalDeaths),
			MapsPlayed:     stat.OverallMaps,
		})
	}

	var hpK, hpD, sndK, sndD, ctlK, ctlD int
	for _, sp := range splits {
		switch sp.Mode {
		case "hp":
			hpK, hpD = sp.Kills, sp.Deaths
		case "snd":
			sndK, sndD = sp.Kills, sp.Deaths
		case "control":
			ctlK, ctlD = sp.Kills, sp.Deaths
		}
	}

	return &PlayerKDStats{
		PlayerID:       playerID,
		Gamertag:       player.Gamertag,
		AvatarURL:      player.AvatarURL,
		TotalKills:     totalKills,
		TotalDeaths:    totalDeaths,
		TotalAssists:   totalAssists,
		AvgKD:          CalculateKD(totalKills, totalDeaths),
		HpKDRatio:      CalculateKD(hpK, hpD),
		SndKDRatio:     CalculateKD(sndK, sndD),
		ControlKDRatio: CalculateKD(ctlK, ctlD),
		Tournaments:    entries,
	}, nil
}

func (ps *PlayerService) GetMatchHistory(ctx context.Context, playerID int) (*PlayerMatchHistory, error) {
	matchStats, err := ps.store.ListMatchHistoryRows(ctx, playerID)
	if err != nil {
		return nil, err
	}

	eventsMap := map[uint]*MatchEvent{}
	var eventOrder []uint
	for _, stat := range matchStats {
		match := stat.Match
		tid := match.TournamentID
		if eventsMap[tid] == nil {
			eventsMap[tid] = &MatchEvent{
				Event:        match.Tournament.Name,
				Year:         match.Tournament.StartDate.Year(),
				TournamentID: tid,
				Matches:      []MatchResult{},
			}
			eventOrder = append(eventOrder, tid)
		}

		var opponent, opponentAbbr, result string
		if stat.TeamID == match.Team1ID {
			opponent = match.Team2.Name
			opponentAbbr = match.Team2.Abbreviation
			if match.Team1Score > match.Team2Score {
				result = "W"
			} else {
				result = "L"
			}
		} else {
			opponent = match.Team1.Name
			opponentAbbr = match.Team1.Abbreviation
			if match.Team2Score > match.Team1Score {
				result = "W"
			} else {
				result = "L"
			}
		}
		var myScore, oppScore int
		if stat.TeamID == match.Team1ID {
			myScore, oppScore = match.Team1Score, match.Team2Score
		} else {
			myScore, oppScore = match.Team2Score, match.Team1Score
		}
		resultScore := result + " " + strconv.Itoa(myScore) + ":" + strconv.Itoa(oppScore)

		eventsMap[tid].Matches = append(eventsMap[tid].Matches, MatchResult{
			MatchID:      match.ID,
			Date:         match.MatchDate.Format(time.RFC3339),
			Opponent:     opponent,
			OpponentAbbr: opponentAbbr,
			Result:       resultScore,
			KD:           stat.KDRatio,
			Kills:        stat.TotalKills,
			Deaths:       stat.TotalDeaths,
		})
	}

	events := make([]MatchEvent, 0, len(eventOrder))
	for _, tid := range eventOrder {
		events = append(events, *eventsMap[tid])
	}
	sortEventsByDate(events)

	return &PlayerMatchHistory{
		PlayerID: playerID,
		Events:   events,
		Total:    len(matchStats),
	}, nil
}

func sortEventsByDate(events []MatchEvent) {
	sort.Slice(events, func(i, j int) bool {
		if len(events[i].Matches) == 0 {
			return false
		}
		if len(events[j].Matches) == 0 {
			return true
		}
		return events[i].Matches[0].Date > events[j].Matches[0].Date
	})
}

// PlayerCareerResult is the assembled franchise-career response for a player.
type PlayerCareerResult struct {
	PlayerID   int                    `json:"player_id"`
	Gamertag   string                 `json:"gamertag"`
	Franchises []FranchiseCareerEntry `json:"franchises"`
}

type FranchiseCareerEntry struct {
	FranchiseKey  string     `json:"franchise_key"`
	FranchiseName string     `json:"franchise_name"`
	Eras          []EraStats `json:"eras"`
	TotalMatches  int        `json:"total_matches"`
	TotalMaps     int        `json:"total_maps"`
	TotalKills    int        `json:"total_kills"`
	TotalDeaths   int        `json:"total_deaths"`
	CareerKD      float64    `json:"career_kd"`
}

type EraStats struct {
	TeamID     uint    `json:"team_id"`
	TeamName   string  `json:"team_name"`
	GameCode   string  `json:"game_code"`
	SeasonName string  `json:"season_name"`
	Matches    int     `json:"matches"`
	Maps       int     `json:"maps"`
	Kills      int     `json:"kills"`
	Deaths     int     `json:"deaths"`
	KD         float64 `json:"kd"`
}

func (ps *PlayerService) GetFranchiseCareer(ctx context.Context, playerID int) (*PlayerCareerResult, error) {
	player, err := ps.store.GetByID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	rows, err := ps.store.ListCareerRows(ctx, playerID)
	if err != nil {
		return nil, err
	}

	franchiseMap := map[string]*FranchiseCareerEntry{}
	var franchiseOrder []string

	for _, r := range rows {
		key := r.FranchiseKey
		if key == "" {
			key = "misc"
		}
		if _, ok := franchiseMap[key]; !ok {
			name := r.FranchiseName
			if key == "misc" {
				name = "Non-CDL / Other"
			}
			franchiseMap[key] = &FranchiseCareerEntry{
				FranchiseKey:  key,
				FranchiseName: name,
				Eras:          []EraStats{},
			}
			franchiseOrder = append(franchiseOrder, key)
		}
		franchiseMap[key].Eras = append(franchiseMap[key].Eras, EraStats{
			TeamID:     r.TeamID,
			TeamName:   r.TeamName,
			GameCode:   r.GameCode,
			SeasonName: r.SeasonName,
			Matches:    r.Matches,
			Maps:       r.Maps,
			Kills:      r.Kills,
			Deaths:     r.Deaths,
			KD:         CalculateKD(r.Kills, r.Deaths),
		})
		franchiseMap[key].TotalMatches += r.Matches
		franchiseMap[key].TotalMaps += r.Maps
		franchiseMap[key].TotalKills += r.Kills
		franchiseMap[key].TotalDeaths += r.Deaths
	}

	result := make([]FranchiseCareerEntry, 0, len(franchiseOrder))
	for _, key := range franchiseOrder {
		f := franchiseMap[key]
		f.CareerKD = CalculateKD(f.TotalKills, f.TotalDeaths)
		result = append(result, *f)
	}

	return &PlayerCareerResult{
		PlayerID:   playerID,
		Gamertag:   player.Gamertag,
		Franchises: result,
	}, nil
}
