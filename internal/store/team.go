package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type TeamStore interface {
	ListActiveCDL(ctx context.Context) ([]models.Team, error)
	ListForSeason(ctx context.Context, seasonID, scope string) ([]models.Team, error)
	ListAll(ctx context.Context) ([]models.Team, error)
	GetByID(ctx context.Context, id int) (*models.Team, error)
	GetPlayers(ctx context.Context, teamID int, seasonID string) ([]models.Player, error)
	GetStats(ctx context.Context, teamID int) ([]models.TeamTournamentStats, error)
}

type gormTeamStore struct{ db *gorm.DB }

func NewGormTeamStore(db *gorm.DB) TeamStore { return &gormTeamStore{db: db} }

func (s *gormTeamStore) ListActiveCDL(ctx context.Context) ([]models.Team, error) {
	var teams []models.Team
	err := s.db.WithContext(ctx).Raw(`
		SELECT DISTINCT t.*
		FROM teams t
		WHERE t.is_cdl_franchise = true
		  AND t.id IN (
		    SELECT DISTINCT pms.team_id
		    FROM player_match_stats pms
		    JOIN matches m ON m.id = pms.match_id
		    JOIN tournaments trn ON trn.id = m.tournament_id
		    JOIN seasons s ON s.id = trn.season_id
		    WHERE s.is_active = true
		      AND trn.tournament_type IN (
		        'major_tournament','qualifier','championship',
		        'kickoff','minor_tournament'
		      )
		  )
		ORDER BY t.name ASC
	`).Scan(&teams).Error
	return teams, err
}

func (s *gormTeamStore) ListForSeason(ctx context.Context, seasonID, scope string) ([]models.Team, error) {
	var sql string
	if scope == "era" {
		sql = `
			SELECT DISTINCT t.*
			FROM teams t
			WHERE t.is_cdl_franchise = true
			  AND t.id IN (
			    SELECT DISTINCT pms.team_id
			    FROM player_match_stats pms
			    JOIN matches m ON m.id = pms.match_id
			    JOIN tournaments trn ON trn.id = m.tournament_id
			    WHERE trn.season_id = ?
			      AND trn.tournament_type NOT IN ('season_summary','unknown')
			  )
			ORDER BY t.name ASC`
	} else {
		sql = `
			SELECT DISTINCT t.*
			FROM teams t
			WHERE t.is_cdl_franchise = true
			  AND t.id IN (
			    SELECT DISTINCT pms.team_id
			    FROM player_match_stats pms
			    JOIN matches m ON m.id = pms.match_id
			    JOIN tournaments trn ON trn.id = m.tournament_id
			    WHERE trn.season_id = ?
			      AND trn.tournament_type IN (
			        'major_tournament','qualifier','championship',
			        'kickoff','minor_tournament'
			      )
			  )
			ORDER BY t.name ASC`
	}
	var teams []models.Team
	err := s.db.WithContext(ctx).Raw(sql, seasonID).Scan(&teams).Error
	return teams, err
}

func (s *gormTeamStore) ListAll(ctx context.Context) ([]models.Team, error) {
	var teams []models.Team
	err := s.db.WithContext(ctx).Order("name ASC").Find(&teams).Error
	return teams, err
}

func (s *gormTeamStore) GetByID(ctx context.Context, id int) (*models.Team, error) {
	var team models.Team
	if err := s.db.WithContext(ctx).Preload("Franchise").First(&team, id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (s *gormTeamStore) GetPlayers(ctx context.Context, teamID int, seasonID string) ([]models.Player, error) {
	query := s.db.WithContext(ctx).
		Joins("JOIN team_rosters ON players.id = team_rosters.player_id").
		Where("team_rosters.team_id = ?", teamID)
	if seasonID != "" {
		query = query.Where("team_rosters.season_id = ?", seasonID)
	} else {
		query = query.Where("team_rosters.end_date IS NULL")
	}
	var players []models.Player
	err := query.Find(&players).Error
	return players, err
}

func (s *gormTeamStore) GetStats(ctx context.Context, teamID int) ([]models.TeamTournamentStats, error) {
	var stats []models.TeamTournamentStats
	err := s.db.WithContext(ctx).
		Where("team_id = ?", teamID).
		Preload("Tournament").
		Find(&stats).Error
	return stats, err
}
