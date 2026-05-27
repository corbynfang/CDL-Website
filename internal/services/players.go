package services

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

// CalculateKD returns kills/deaths, or 0 if deaths is 0.
func CalculateKD(kills, deaths int) float64 {
	if deaths == 0 {
		return 0
	}
	return float64(kills) / float64(deaths)
}

func calculateKD(kills, deaths int) float64 { return CalculateKD(kills, deaths) }

type PlayerService struct {
	store store.PlayerStore
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

// PlayerKDStats is the assembled K/D profile returned by GetKDStats.
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

func (ps *PlayerService) GetKDStats(ctx context.Context, playerID int) (*PlayerKDStats, error) {
	player, err := ps.store.GetByID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	tournamentStats, err := ps.store.ListTournamentStats(ctx, playerID)
	if err != nil {
		return nil, err
	}

	var totalKills, totalDeaths, totalAssists int
	var totalHpK, totalHpD, totalSndK, totalSndD int
	var ctlKDSum float64
	var ctlMapsTotal int

	entries := make([]TournamentKDEntry, 0, len(tournamentStats))
	for _, stat := range tournamentStats {
		totalKills += stat.TotalKills
		totalDeaths += stat.TotalDeaths
		totalAssists += stat.TotalAssists
		totalHpK += stat.HpKills
		totalHpD += stat.HpDeaths
		totalSndK += stat.SndKills
		totalSndD += stat.SndDeaths
		if stat.ControlMaps > 0 {
			ctlKDSum += stat.ControlKDRatio * float64(stat.ControlMaps)
			ctlMapsTotal += stat.ControlMaps
		}
		entries = append(entries, TournamentKDEntry{
			TournamentID:   stat.TournamentID,
			TournamentName: stat.Tournament.Name,
			Kills:          stat.TotalKills,
			Deaths:         stat.TotalDeaths,
			Assists:        stat.TotalAssists,
			KDRatio:        calculateKD(stat.TotalKills, stat.TotalDeaths),
			MapsPlayed:     stat.OverallMaps,
		})
	}

	controlKD := 0.0
	if ctlMapsTotal > 0 {
		controlKD = ctlKDSum / float64(ctlMapsTotal)
	}

	return &PlayerKDStats{
		PlayerID:       playerID,
		Gamertag:       player.Gamertag,
		AvatarURL:      player.AvatarURL,
		TotalKills:     totalKills,
		TotalDeaths:    totalDeaths,
		TotalAssists:   totalAssists,
		AvgKD:          calculateKD(totalKills, totalDeaths),
		HpKDRatio:      calculateKD(totalHpK, totalHpD),
		SndKDRatio:     calculateKD(totalSndK, totalSndD),
		ControlKDRatio: controlKD,
		Tournaments:    entries,
	}, nil
}

// PlayerMatchHistory is the assembled match-history response for a single player.
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
		resultScore := result + " " + strconv.Itoa(match.Team1Score) + ":" + strconv.Itoa(match.Team2Score)

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
			KD:         calculateKD(r.Kills, r.Deaths),
		})
		franchiseMap[key].TotalMatches += r.Matches
		franchiseMap[key].TotalMaps += r.Maps
		franchiseMap[key].TotalKills += r.Kills
		franchiseMap[key].TotalDeaths += r.Deaths
	}

	result := make([]FranchiseCareerEntry, 0, len(franchiseOrder))
	for _, key := range franchiseOrder {
		f := franchiseMap[key]
		f.CareerKD = calculateKD(f.TotalKills, f.TotalDeaths)
		result = append(result, *f)
	}

	return &PlayerCareerResult{
		PlayerID:   playerID,
		Gamertag:   player.Gamertag,
		Franchises: result,
	}, nil
}
