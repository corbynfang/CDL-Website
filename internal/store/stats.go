package store

import (
	"context"

	"gorm.io/gorm"
)

type KDRow struct {
	PlayerID      uint    `json:"player_id"`
	Gamertag      string  `json:"gamertag"`
	AvatarURL     string  `json:"avatar_url"`
	TeamAbbr      string  `json:"team_abbr"`
	SeasonKills   int     `json:"season_kills"`
	SeasonDeaths  int     `json:"season_deaths"`
	SeasonAssists int     `json:"season_assists"`
	SeasonKD      float64 `json:"season_kd"`
}

type StatsStore interface {
	GetTopKDRows(ctx context.Context, limit int) ([]KDRow, error)
	GetAllKDRows(ctx context.Context, limit int, seasonID string) ([]KDRow, error)
}

type gormStatsStore struct{ db *gorm.DB }

func NewGormStatsStore(db *gorm.DB) StatsStore { return &gormStatsStore{db: db} }

func (s *gormStatsStore) kdBase(ctx context.Context) *gorm.DB {
	return s.db.WithContext(ctx).
		Table("player_tournament_stats pts").
		Select(`pts.player_id, MAX(p.gamertag) as gamertag,
			COALESCE(MAX(p.avatar_url), '') as avatar_url,
			COALESCE(MAX(t.abbreviation), '') as team_abbr,
			SUM(pts.total_kills) as season_kills,
			SUM(pts.total_deaths) as season_deaths,
			SUM(pts.total_assists) as season_assists`).
		Joins("JOIN players p ON pts.player_id = p.id").
		Joins("LEFT JOIN teams t ON pts.team_id = t.id").
		Group("pts.player_id")
}

func (s *gormStatsStore) GetTopKDRows(ctx context.Context, limit int) ([]KDRow, error) {
	rows := make([]KDRow, 0)
	err := s.kdBase(ctx).
		Having("SUM(pts.total_deaths) > 0").
		Order("(SUM(pts.total_kills)::decimal / NULLIF(SUM(pts.total_deaths), 0)) DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (s *gormStatsStore) GetAllKDRows(ctx context.Context, limit int, seasonID string) ([]KDRow, error) {
	query := s.kdBase(ctx).
		Joins("JOIN tournaments tour ON pts.tournament_id = tour.id")
	if seasonID != "" {
		query = query.Where("tour.season_id = ?", seasonID)
	}
	rows := make([]KDRow, 0)
	err := query.
		Group("pts.player_id").
		Having("SUM(pts.total_kills) > 0 OR SUM(pts.total_deaths) > 0").
		Order("(CASE WHEN SUM(pts.total_deaths) > 0 THEN SUM(pts.total_kills)::decimal / SUM(pts.total_deaths) ELSE 0 END) DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}
