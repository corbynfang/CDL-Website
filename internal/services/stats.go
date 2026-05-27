package services

import (
	"context"

	"github.com/corbynfang/CDL-Website/internal/store"
)

// PlayerKDRow is a type alias for store.KDRow — the same struct, same JSON tags.
// SeasonKD is computed here after the raw DB scan.
type PlayerKDRow = store.KDRow

type StatsService struct {
	store store.StatsStore
}

func NewStatsService(s store.StatsStore) *StatsService {
	return &StatsService{store: s}
}

func (ss *StatsService) GetTopKD(ctx context.Context, limit int) ([]PlayerKDRow, error) {
	rows, err := ss.store.GetTopKDRows(ctx, limit)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].SeasonKD = calculateKD(rows[i].SeasonKills, rows[i].SeasonDeaths)
	}
	return rows, nil
}

func (ss *StatsService) GetAllKD(ctx context.Context, limit int, seasonID string) ([]PlayerKDRow, error) {
	rows, err := ss.store.GetAllKDRows(ctx, limit, seasonID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].SeasonKD = calculateKD(rows[i].SeasonKills, rows[i].SeasonDeaths)
	}
	return rows, nil
}

func (ss *StatsService) GetTableCounts(ctx context.Context) (store.TableCounts, error) {
	return ss.store.GetTableCounts(ctx)
}
