package store

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

// MatchStore covers per-match DB operations: the match itself, its maps, and per-map player stats.
type MatchStore interface {
	GetByID(ctx context.Context, id int) (*models.Match, error)
	GetMaps(ctx context.Context, matchID int) ([]models.MatchMap, error)
	GetStatRows(ctx context.Context, matchID int) ([]MatchStatRow, error)
}

// MatchStatRow is the raw scan target for the per-map player-stat query.
type MatchStatRow struct {
	MapNumber       int
	PlayerID        uint
	Gamertag        string
	TeamID          uint
	Kills           int
	Deaths          int
	KDRatio         float64
	Damage          int
	Assists         int
	BPRating        float64
	HillTime        int
	SndRounds       int
	PlantCount      int
	DefuseCount     int
	FirstBloodCount int
	FirstDeathCount int
	NonTradedKills  int
	HighestStreak   int
	DataQualityNote string
}

type gormMatchStore struct{ db *gorm.DB }

func NewGormMatchStore(db *gorm.DB) MatchStore { return &gormMatchStore{db: db} }

func (s *gormMatchStore) GetByID(ctx context.Context, id int) (*models.Match, error) {
	var match models.Match
	err := s.db.WithContext(ctx).
		Preload("Team1").
		Preload("Team2").
		Preload("Winner").
		Preload("Tournament").
		Preload("Tournament.Season").
		First(&match, id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (s *gormMatchStore) GetMaps(ctx context.Context, matchID int) ([]models.MatchMap, error) {
	var maps []models.MatchMap
	err := s.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Order("map_number ASC").
		Find(&maps).Error
	return maps, err
}

func (s *gormMatchStore) GetStatRows(ctx context.Context, matchID int) ([]MatchStatRow, error) {
	var rows []MatchStatRow
	err := s.db.WithContext(ctx).
		Table("player_map_stats pms").
		Select(`pms.map_number, pms.player_id, p.gamertag, pms.team_id,
			pms.kills, pms.deaths, pms.kd_ratio, pms.damage, pms.assists,
			pms.bp_rating, pms.hill_time, pms.snd_rounds, pms.plant_count,
			pms.defuse_count, pms.first_blood_count, pms.first_death_count,
			pms.non_traded_kills, pms.highest_streak, pms.data_quality_note`).
		Joins("JOIN players p ON p.id = pms.player_id").
		Where("pms.match_id = ?", matchID).
		Order("pms.map_number ASC, pms.kills DESC").
		Scan(&rows).Error
	return rows, err
}
