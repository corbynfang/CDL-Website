package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

// TournamentStore covers tournament-level DB operations used by the MatchService.
type TournamentStore interface {
	List(ctx context.Context, seasonID string) ([]models.Tournament, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tournament, error)
	GetByID(ctx context.Context, id int) (*models.Tournament, error)
	GetTeamCount(ctx context.Context, tournamentID int) (int64, error)
	GetBracketMatches(ctx context.Context, tournamentID int) ([]models.Match, error)
	GetMatches(ctx context.Context, tournamentID int) ([]models.Match, error)
	GetTeamIDs(ctx context.Context, tournamentID int) ([]uint, error)
	GetTeams(ctx context.Context, teamIDs []uint) ([]models.Team, error)
	GetTeamStats(ctx context.Context, tournamentID int) ([]models.TeamTournamentStats, error)
	GetPlayerStats(ctx context.Context, tournamentID int) ([]models.PlayerTournamentStats, error)
}

type gormTournamentStore struct{ db *gorm.DB }

func NewGormTournamentStore(db *gorm.DB) TournamentStore { return &gormTournamentStore{db: db} }

func (s *gormTournamentStore) List(ctx context.Context, seasonID string) ([]models.Tournament, error) {
	query := s.db.WithContext(ctx).
		Preload("Season").
		Where("tournament_type NOT IN ('season_summary','unknown')").
		Order("start_date DESC")
	if seasonID != "" {
		query = query.Where("season_id = ?", seasonID)
	}
	var tournaments []models.Tournament
	err := query.Find(&tournaments).Error
	return tournaments, err
}

func (s *gormTournamentStore) GetBySlug(ctx context.Context, slug string) (*models.Tournament, error) {
	var tournament models.Tournament
	err := s.db.WithContext(ctx).
		Preload("Season").
		Where("slug = ?", slug).
		First(&tournament).Error
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (s *gormTournamentStore) GetByID(ctx context.Context, id int) (*models.Tournament, error) {
	var tournament models.Tournament
	if err := s.db.WithContext(ctx).Preload("Season").First(&tournament, id).Error; err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (s *gormTournamentStore) GetTeamCount(ctx context.Context, tournamentID int) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM (
			SELECT team1_id AS team_id FROM matches WHERE tournament_id = ?
			UNION
			SELECT team2_id FROM matches WHERE tournament_id = ?
		) AS t
	`, tournamentID, tournamentID).Scan(&count).Error
	return count, err
}

func (s *gormTournamentStore) GetBracketMatches(ctx context.Context, tournamentID int) ([]models.Match, error) {
	var matches []models.Match
	err := s.db.WithContext(ctx).
		Where("tournament_id = ? AND bracket_round != ''", tournamentID).
		Preload("Team1").
		Preload("Team2").
		Order("bracket_round, bracket_position").
		Find(&matches).Error
	return matches, err
}

func (s *gormTournamentStore) GetMatches(ctx context.Context, tournamentID int) ([]models.Match, error) {
	var matches []models.Match
	err := s.db.WithContext(ctx).
		Where("tournament_id = ?", tournamentID).
		Preload("Team1").
		Preload("Team2").
		Preload("Winner").
		Order("match_date ASC, bracket_position ASC").
		Find(&matches).Error
	return matches, err
}

func (s *gormTournamentStore) GetTeamIDs(ctx context.Context, tournamentID int) ([]uint, error) {
	var teamIDs []uint
	err := s.db.WithContext(ctx).Raw(`
		SELECT DISTINCT team_id FROM (
			SELECT team1_id AS team_id FROM matches WHERE tournament_id = ?
			UNION ALL
			SELECT team2_id FROM matches WHERE tournament_id = ?
		) AS t
	`, tournamentID, tournamentID).Scan(&teamIDs).Error
	return teamIDs, err
}

func (s *gormTournamentStore) GetTeams(ctx context.Context, teamIDs []uint) ([]models.Team, error) {
	var teams []models.Team
	err := s.db.WithContext(ctx).Where("id IN ?", teamIDs).Find(&teams).Error
	return teams, err
}

func (s *gormTournamentStore) GetTeamStats(ctx context.Context, tournamentID int) ([]models.TeamTournamentStats, error) {
	var stats []models.TeamTournamentStats
	err := s.db.WithContext(ctx).Where("tournament_id = ?", tournamentID).Find(&stats).Error
	return stats, err
}

func (s *gormTournamentStore) GetPlayerStats(ctx context.Context, tournamentID int) ([]models.PlayerTournamentStats, error) {
	var stats []models.PlayerTournamentStats
	err := s.db.WithContext(ctx).
		Where("tournament_id = ? AND (total_kills > 0 OR total_deaths > 0)", tournamentID).
		Preload("Player").
		Preload("Team").
		Order("(CASE WHEN total_deaths > 0 THEN CAST(total_kills AS decimal) / total_deaths ELSE 0 END) DESC").
		Find(&stats).Error
	return stats, err
}
