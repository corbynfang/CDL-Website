package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/corbynfang/CDL-Website/internal/store"
)

// PlayerKDRow is a type alias for store.KDRow — the same struct, same JSON tags.
// SeasonKD is computed here after the raw DB scan.
type PlayerKDRow = store.KDRow

const statsCacheTTL = 5 * time.Minute

type kdCacheEntry struct {
	rows []PlayerKDRow
	exp  time.Time
}

type StatsService struct {
	store store.StatsStore
	mu    sync.RWMutex
	topKD map[int]kdCacheEntry
	allKD map[string]kdCacheEntry
}

func NewStatsService(s store.StatsStore) *StatsService {
	return &StatsService{
		store: s,
		topKD: make(map[int]kdCacheEntry),
		allKD: make(map[string]kdCacheEntry),
	}
}

func (ss *StatsService) GetTopKD(ctx context.Context, limit int) ([]PlayerKDRow, error) {
	ss.mu.RLock()
	if entry, ok := ss.topKD[limit]; ok && time.Now().Before(entry.exp) {
		rows := entry.rows
		ss.mu.RUnlock()
		return rows, nil
	}
	ss.mu.RUnlock()

	rows, err := ss.store.GetTopKDRows(ctx, limit)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].SeasonKD = CalculateKD(rows[i].SeasonKills, rows[i].SeasonDeaths)
	}

	ss.mu.Lock()
	ss.topKD[limit] = kdCacheEntry{rows: rows, exp: time.Now().Add(statsCacheTTL)}
	ss.mu.Unlock()

	return rows, nil
}

func (ss *StatsService) GetAllKD(ctx context.Context, limit int, seasonID string) ([]PlayerKDRow, error) {
	key := fmt.Sprintf("%d:%s", limit, seasonID)

	ss.mu.RLock()
	if entry, ok := ss.allKD[key]; ok && time.Now().Before(entry.exp) {
		rows := entry.rows
		ss.mu.RUnlock()
		return rows, nil
	}
	ss.mu.RUnlock()

	rows, err := ss.store.GetAllKDRows(ctx, limit, seasonID)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].SeasonKD = CalculateKD(rows[i].SeasonKills, rows[i].SeasonDeaths)
	}

	ss.mu.Lock()
	ss.allKD[key] = kdCacheEntry{rows: rows, exp: time.Now().Add(statsCacheTTL)}
	ss.mu.Unlock()

	return rows, nil
}
