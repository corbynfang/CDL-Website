package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

// PlayerStore is the data-access contract the PlayerService depends on.
// Every method maps to one distinct DB operation; business logic lives in the service.
type PlayerStore interface {
	List(ctx context.Context, search string, limit, offset int) ([]models.Player, int64, error)
	GetByID(ctx context.Context, id int) (*models.Player, error)
	ListMatchStats(ctx context.Context, playerID int) ([]models.PlayerMatchStats, error)
	ListTournamentStats(ctx context.Context, playerID int) ([]models.PlayerTournamentStats, error)
	ListModeKDSplits(ctx context.Context, playerID int) ([]ModeKDSplit, error)
	ListMatchHistoryRows(ctx context.Context, playerID int) ([]models.PlayerMatchStats, error)
	ListCareerRows(ctx context.Context, playerID int) ([]PlayerCareerRow, error)
}

// ModeKDSplit holds a player's kill/death totals for one game mode, aggregated
// from per-map stats. Mode is normalised to "hp", "snd", or "control".
type ModeKDSplit struct {
	Mode   string
	Kills  int
	Deaths int
}

// PlayerCareerRow is the raw scan target for the franchise-career SQL query.
// Business-logic aggregation (totals, KD) happens in the service after this is returned.
type PlayerCareerRow struct {
	FranchiseID   *uint
	FranchiseKey  string
	FranchiseName string
	TeamID        uint
	TeamName      string
	GameCode      string
	SeasonName    string
	Matches       int
	Maps          int
	Kills         int
	Deaths        int
}

type gormPlayerStore struct{ db *gorm.DB }

func NewGormPlayerStore(db *gorm.DB) PlayerStore { return &gormPlayerStore{db: db} }

func (s *gormPlayerStore) List(ctx context.Context, search string, limit, offset int) ([]models.Player, int64, error) {
	base := s.db.WithContext(ctx).Model(&models.Player{})
	if search != "" {
		base = base.Where("gamertag ILIKE ?", "%"+search+"%")
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var players []models.Player
	err := base.Order("gamertag ASC").Limit(limit).Offset(offset).Find(&players).Error
	return players, total, err
}

func (s *gormPlayerStore) GetByID(ctx context.Context, id int) (*models.Player, error) {
	var player models.Player
	if err := s.db.WithContext(ctx).First(&player, id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (s *gormPlayerStore) ListMatchStats(ctx context.Context, playerID int) ([]models.PlayerMatchStats, error) {
	var stats []models.PlayerMatchStats
	err := s.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("Match").
		Preload("Team").
		Find(&stats).Error
	return stats, err
}

func (s *gormPlayerStore) ListTournamentStats(ctx context.Context, playerID int) ([]models.PlayerTournamentStats, error) {
	var stats []models.PlayerTournamentStats
	err := s.db.WithContext(ctx).
		Where("player_id = ?", playerID).
		Preload("Tournament").
		Order("tournament_id DESC").
		Find(&stats).Error
	return stats, err
}

func (s *gormPlayerStore) ListModeKDSplits(ctx context.Context, playerID int) ([]ModeKDSplit, error) {
	var rows []ModeKDSplit
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			CASE
				WHEN mm.mode IN ('Search and Destroy', 'Search & Destroy') THEN 'snd'
				WHEN mm.mode = 'Hardpoint' THEN 'hp'
				WHEN mm.mode = 'Control'   THEN 'control'
				ELSE 'other'
			END AS mode,
			COALESCE(SUM(pms.kills), 0)  AS kills,
			COALESCE(SUM(pms.deaths), 0) AS deaths
		FROM player_map_stats pms
		JOIN match_maps mm
			ON mm.match_id = pms.match_id
			AND mm.map_number = pms.map_number
		WHERE pms.player_id = ?
			AND mm.mode IS NOT NULL
			AND mm.mode <> ''
			AND mm.played = true
		GROUP BY 1
	`, playerID).Scan(&rows).Error
	return rows, err
}

func (s *gormPlayerStore) ListMatchHistoryRows(ctx context.Context, playerID int) ([]models.PlayerMatchStats, error) {
	var stats []models.PlayerMatchStats
	err := s.db.WithContext(ctx).
		Where("player_match_stats.player_id = ?", playerID).
		Preload("Match").
		Preload("Match.Tournament").
		Preload("Match.Team1").
		Preload("Match.Team2").
		Preload("Team").
		Joins("JOIN matches ON matches.id = player_match_stats.match_id").
		Order("CASE WHEN matches.match_date <= '0001-01-02 00:00:00+00'::timestamptz THEN 0 ELSE 1 END DESC, matches.match_date DESC, player_match_stats.match_id DESC").
		Limit(500).
		Find(&stats).Error
	return stats, err
}

func (s *gormPlayerStore) ListCareerRows(ctx context.Context, playerID int) ([]PlayerCareerRow, error) {
	var rows []PlayerCareerRow
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			t.franchise_id,
			COALESCE(f.franchise_key, '')      AS franchise_key,
			COALESCE(f.name, t.name)           AS franchise_name,
			pms.team_id                        AS team_id,
			t.name                             AS team_name,
			COALESCE(t.game_code, '')          AS game_code,
			COALESCE(s.name, '')               AS season_name,
			COUNT(DISTINCT pms.match_id)       AS matches,
			SUM(pms.maps_played)               AS maps,
			SUM(pms.total_kills)               AS kills,
			SUM(pms.total_deaths)              AS deaths
		FROM player_match_stats pms
		JOIN teams t ON t.id = pms.team_id
		LEFT JOIN franchises f ON f.id = t.franchise_id
		LEFT JOIN matches m ON m.id = pms.match_id
		LEFT JOIN tournaments tour ON tour.id = m.tournament_id
		LEFT JOIN seasons s ON s.id = tour.season_id
		WHERE pms.player_id = ?
		GROUP BY t.franchise_id, f.franchise_key, f.name, pms.team_id, t.name, t.game_code, s.name
		ORDER BY t.franchise_id NULLS LAST, t.game_code
	`, playerID).Scan(&rows).Error
	return rows, err
}
