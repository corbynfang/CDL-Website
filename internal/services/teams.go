package services

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/corbynfang/CDL-Website/internal/store"
)

var ErrInvalidSeason = errors.New("invalid season")

type TeamService struct {
	teams   store.TeamStore
	seasons store.SeasonStore

	cacheMu   sync.RWMutex
	cached    []models.Team
	cacheExp  time.Time
}

const defaultTeamsCacheTTL = 60 * time.Second

func NewTeamService(teams store.TeamStore, seasons store.SeasonStore) *TeamService {
	return &TeamService{teams: teams, seasons: seasons}
}

func (ts *TeamService) List(ctx context.Context, seasonID, scope string) ([]models.Team, error) {
	if seasonID != "" {
		id, err := strconv.Atoi(seasonID)
		if err != nil {
			return nil, ErrInvalidSeason
		}
		season, err := ts.seasons.GetByID(ctx, id)
		if err != nil || season.GameCode == "" {
			return nil, ErrInvalidSeason
		}
		teams, err := ts.teams.ListForSeason(ctx, seasonID, scope)
		if err != nil {
			return nil, err
		}
		sort.Slice(teams, func(i, j int) bool { return teams[i].Name < teams[j].Name })
		return teams, nil
	}

	if scope == "all" {
		return ts.teams.ListAll(ctx)
	}

	ts.cacheMu.RLock()
	cached := ts.cached
	valid := time.Now().Before(ts.cacheExp)
	ts.cacheMu.RUnlock()

	if valid && cached != nil {
		return cached, nil
	}

	teams, err := ts.teams.ListActiveCDL(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(teams, func(i, j int) bool { return teams[i].Name < teams[j].Name })

	ts.cacheMu.Lock()
	ts.cached = teams
	ts.cacheExp = time.Now().Add(defaultTeamsCacheTTL)
	ts.cacheMu.Unlock()

	return teams, nil
}

func (ts *TeamService) GetByID(ctx context.Context, id int) (*models.Team, error) {
	return ts.teams.GetByID(ctx, id)
}

func (ts *TeamService) GetPlayers(ctx context.Context, teamID int, seasonID string) ([]models.Player, error) {
	return ts.teams.GetPlayers(ctx, teamID, seasonID)
}

func (ts *TeamService) GetCurrentRoster(ctx context.Context, teamID int, seasonID string) ([]models.Player, error) {
	return ts.teams.GetLatestMatchRoster(ctx, teamID, seasonID)
}

func (ts *TeamService) GetStats(ctx context.Context, teamID int) ([]models.TeamTournamentStats, error) {
	return ts.teams.GetStats(ctx, teamID)
}
