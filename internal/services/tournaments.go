package services

import (
	"context"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

type TournamentDetail struct {
	Tournament  models.Tournament `json:"tournament"`
	TeamCount   int64             `json:"team_count"`
	EventFormat string            `json:"event_format"`
}

type TournamentTeam struct {
	models.Team
	Placement   *int `json:"placement"`
	MatchesWon  int  `json:"matches_won"`
	MatchesLost int  `json:"matches_lost"`
}

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

type TournamentService struct {
	tournaments store.TournamentStore
}

func NewTournamentService(tournaments store.TournamentStore) *TournamentService {
	return &TournamentService{tournaments: tournaments}
}

func (ts *TournamentService) ListTournaments(ctx context.Context, seasonID string) ([]models.Tournament, error) {
	return ts.tournaments.List(ctx, seasonID)
}

func (ts *TournamentService) GetTournamentBySlug(ctx context.Context, slug string) (*TournamentDetail, error) {
	tournament, err := ts.tournaments.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	teamCount, err := ts.tournaments.GetTeamCount(ctx, int(tournament.ID))
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

func (ts *TournamentService) GetTournamentByID(ctx context.Context, id int) (*models.Tournament, error) {
	return ts.tournaments.GetByID(ctx, id)
}

func (ts *TournamentService) AssembleBracket(ctx context.Context, tournamentID int) (*BracketResult, error) {
	tournament, err := ts.tournaments.GetByID(ctx, tournamentID)
	if err != nil {
		return nil, err
	}

	matches, err := ts.tournaments.GetBracketMatches(ctx, tournamentID)
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

func (ts *TournamentService) ListTournamentMatches(ctx context.Context, tournamentID int) ([]models.Match, error) {
	return ts.tournaments.GetMatches(ctx, tournamentID)
}

func (ts *TournamentService) GetTournamentTeams(ctx context.Context, tournamentID int) ([]TournamentTeam, error) {
	teamIDs, err := ts.tournaments.GetTeamIDs(ctx, tournamentID)
	if err != nil {
		return nil, err
	}
	if len(teamIDs) == 0 {
		return []TournamentTeam{}, nil
	}

	teams, err := ts.tournaments.GetTeams(ctx, teamIDs)
	if err != nil {
		return nil, err
	}

	stats, err := ts.tournaments.GetTeamStats(ctx, tournamentID)
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

func (ts *TournamentService) GetTournamentStats(ctx context.Context, tournamentID int) ([]models.PlayerTournamentStats, error) {
	return ts.tournaments.GetPlayerStats(ctx, tournamentID)
}
