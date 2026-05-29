package store

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

type TeamStore interface {
	ListActiveCDL(ctx context.Context) ([]models.Team, error)
	ListForSeason(ctx context.Context, seasonID, scope string) ([]models.Team, error)
	ListAll(ctx context.Context) ([]models.Team, error)
	GetByID(ctx context.Context, id int) (*models.Team, error)
	GetPlayers(ctx context.Context, teamID int, seasonID string) ([]models.Player, error)
	GetLatestMatchRoster(ctx context.Context, teamID int, seasonID string) ([]models.Player, error)
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
		Where("team_rosters.team_id = ?", teamID).
		Distinct().
		Order("players.gamertag ASC")
	if seasonID != "" {
		query = query.Where("team_rosters.season_id = ?", seasonID)
	} else {
		query = query.Where(`team_rosters.season_id = (
			SELECT tr2.season_id
			FROM team_rosters tr2
			JOIN seasons s2 ON s2.id = tr2.season_id
			WHERE tr2.team_id = ?
			ORDER BY s2.start_date DESC
			LIMIT 1
		)`, teamID)
	}
	var players []models.Player
	err := query.Find(&players).Error
	return players, err
}

func (s *gormTeamStore) GetLatestMatchRoster(ctx context.Context, teamID int, seasonID string) ([]models.Player, error) {
	season := seasonID
	if season == "" {
		latest, err := s.latestSeasonByMatch(ctx, teamID)
		if err != nil {
			return nil, err
		}
		season = latest
	}

	if season != "" {
		players, err := s.latestMatchRoster(ctx, teamID, season)
		if err != nil {
			return nil, err
		}
		if len(players) > 0 {
			return players, nil
		}
	}

	return s.GetPlayers(ctx, teamID, seasonID)
}

func (s *gormTeamStore) latestSeasonByMatch(ctx context.Context, teamID int) (string, error) {
	var sid sql.NullInt64
	err := s.db.WithContext(ctx).Raw(`
		SELECT tour.season_id
		FROM player_map_stats pms
		JOIN matches m ON m.id = pms.match_id
		JOIN tournaments tour ON tour.id = m.tournament_id
		WHERE pms.team_id = ?
		ORDER BY m.match_date DESC, m.id DESC
		LIMIT 1
	`, teamID).Scan(&sid).Error
	if err != nil || !sid.Valid {
		return "", err
	}
	return strconv.FormatInt(sid.Int64, 10), nil
}

func (s *gormTeamStore) latestMatchRoster(ctx context.Context, teamID int, seasonID string) ([]models.Player, error) {
	var players []models.Player
	err := s.db.WithContext(ctx).Raw(`
		SELECT DISTINCT p.*
		FROM players p
		JOIN player_map_stats pms ON pms.player_id = p.id
		LEFT JOIN match_maps mm ON mm.match_id = pms.match_id AND mm.map_number = pms.map_number
		WHERE pms.team_id = ?
		  AND (mm.id IS NULL OR mm.played = true)
		  AND pms.match_id = (
		    SELECT m.id
		    FROM player_map_stats pms2
		    JOIN matches m ON m.id = pms2.match_id
		    JOIN tournaments tour ON tour.id = m.tournament_id
		    LEFT JOIN match_maps mm2 ON mm2.match_id = pms2.match_id AND mm2.map_number = pms2.map_number
		    WHERE pms2.team_id = ? AND tour.season_id = ?
		      AND (mm2.id IS NULL OR mm2.played = true)
		    ORDER BY m.match_date DESC, m.id DESC
		    LIMIT 1
		  )
		ORDER BY p.gamertag ASC
	`, teamID, teamID, seasonID).Scan(&players).Error
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
