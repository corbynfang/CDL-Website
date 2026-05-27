package services

import (
	"context"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
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
type TournamentDetail struct {
	Tournament  models.Tournament `json:"tournament"`
	TeamCount   int64             `json:"team_count"`
	EventFormat string            `json:"event_format"`
}

// TournamentTeam is a team enriched with its placement and record for a specific tournament.
type TournamentTeam struct {
	models.Team
	Placement   *int `json:"placement"`
	MatchesWon  int  `json:"matches_won"`
	MatchesLost int  `json:"matches_lost"`
}

// BracketResult is the assembled bracket response including optional group stage.
type BracketResult struct {
	TournamentID   int                       `json:"tournament_id"`
	TournamentName string                    `json:"tournament_name"`
	EventFormat    string                    `json:"event_format"`
	TotalMatches   int                       `json:"total_matches"`
	Bracket        map[string][]BracketMatch `json:"bracket"`
	GroupStage     map[string][]BracketMatch `json:"group_stage,omitempty"`
}

type BracketMatch struct {
	ID              uint      `json:"id"`
	Team1ID         uint      `json:"team1_id"`
	Team2ID         uint      `json:"team2_id"`
	Team1Name       string    `json:"team1_name"`
	Team1Abbr       string    `json:"team1_abbr"`
	Team1Logo       string    `json:"team1_logo"`
	Team2Name       string    `json:"team2_name"`
	Team2Abbr       string    `json:"team2_abbr"`
	Team2Logo       string    `json:"team2_logo"`
	Team1Score      int       `json:"team1_score"`
	Team2Score      int       `json:"team2_score"`
	WinnerID        *uint     `json:"winner_id"`
	BracketPosition int       `json:"bracket_position"`
	MatchDate       time.Time `json:"match_date"`
}

type MatchService struct {
	matches     store.MatchStore
	tournaments store.TournamentStore
}

func NewMatchService(matches store.MatchStore, tournaments store.TournamentStore) *MatchService {
	return &MatchService{matches: matches, tournaments: tournaments}
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

func (ms *MatchService) ListTournaments(ctx context.Context, seasonID string) ([]models.Tournament, error) {
	return ms.tournaments.List(ctx, seasonID)
}

func (ms *MatchService) GetTournamentBySlug(ctx context.Context, slug string) (*TournamentDetail, error) {
	tournament, err := ms.tournaments.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	teamCount, err := ms.tournaments.GetTeamCount(ctx, int(tournament.ID))
	if err != nil {
		return nil, err
	}

	if tournament.TournamentFormat == "" {
		if f := detectBracketFormat("", tournament.TournamentType); f != bracketFmtUnknown {
			tournament.TournamentFormat = FormatName(f)
		}
	}

	return &TournamentDetail{
		Tournament:  *tournament,
		TeamCount:   teamCount,
		EventFormat: tournament.TournamentFormat,
	}, nil
}

func (ms *MatchService) GetTournamentByID(ctx context.Context, id int) (*models.Tournament, error) {
	return ms.tournaments.GetByID(ctx, id)
}

func (ms *MatchService) AssembleBracket(ctx context.Context, tournamentID int) (*BracketResult, error) {
	tournament, err := ms.tournaments.GetByID(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	matches, err := ms.tournaments.GetBracketMatches(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	format := detectBracketFormat(tournament.TournamentFormat, tournament.TournamentType)
	normalize := roundNormalizerFor(format)

	bracketKeySet := bracketKeysFor(format)
	bracket := make(map[string][]BracketMatch, len(bracketKeySet))
	for k := range bracketKeySet {
		bracket[k] = []BracketMatch{}
	}

	var groupStage map[string][]BracketMatch
	if hasGroupStage(format) {
		groupStage = map[string][]BracketMatch{}
	}

	for _, match := range matches {
		bm := BracketMatch{
			ID:              match.ID,
			Team1ID:         match.Team1ID,
			Team2ID:         match.Team2ID,
			Team1Name:       match.Team1.Name,
			Team1Abbr:       match.Team1.Abbreviation,
			Team1Logo:       match.Team1.LogoURL,
			Team2Name:       match.Team2.Name,
			Team2Abbr:       match.Team2.Abbreviation,
			Team2Logo:       match.Team2.LogoURL,
			Team1Score:      match.Team1Score,
			Team2Score:      match.Team2Score,
			WinnerID:        match.WinnerID,
			BracketPosition: match.BracketPosition,
			MatchDate:       match.MatchDate,
		}
		key := normalize(match.BracketRound)
		if format == bracketFmtEWCGroupBracket {
			key = ewcGroupKey(key, match.BracketPosition)
		}
		if _, inBracket := bracket[key]; inBracket {
			bracket[key] = append(bracket[key], bm)
		} else if groupStage != nil {
			groupStage[key] = append(groupStage[key], bm)
		}
	}

	return &BracketResult{
		TournamentID:   tournamentID,
		TournamentName: tournament.Name,
		EventFormat:    FormatName(format),
		TotalMatches:   len(matches),
		Bracket:        bracket,
		GroupStage:     groupStage,
	}, nil
}

func (ms *MatchService) ListTournamentMatches(ctx context.Context, tournamentID int) ([]models.Match, error) {
	return ms.tournaments.GetMatches(ctx, tournamentID)
}

func (ms *MatchService) GetTournamentTeams(ctx context.Context, tournamentID int) ([]TournamentTeam, error) {
	teamIDs, err := ms.tournaments.GetTeamIDs(ctx, tournamentID)
	if err != nil {
		return nil, err
	}
	if len(teamIDs) == 0 {
		return []TournamentTeam{}, nil
	}

	teams, err := ms.tournaments.GetTeams(ctx, teamIDs)
	if err != nil {
		return nil, err
	}

	stats, err := ms.tournaments.GetTeamStats(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	statsMap := make(map[uint]models.TeamTournamentStats, len(stats))
	for _, s := range stats {
		statsMap[s.TeamID] = s
	}

	result := make([]TournamentTeam, 0, len(teams))
	for _, t := range teams {
		out := TournamentTeam{Team: t}
		if s, ok := statsMap[t.ID]; ok {
			out.Placement = s.Placement
			out.MatchesWon = s.MatchesWon
			out.MatchesLost = s.MatchesLost
		}
		result = append(result, out)
	}
	return result, nil
}

func (ms *MatchService) GetTournamentStats(ctx context.Context, tournamentID int) ([]models.PlayerTournamentStats, error) {
	return ms.tournaments.GetPlayerStats(ctx, tournamentID)
}
