package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mustUTC(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04", s, time.UTC)
	if err != nil {
		panic("mustUTC: " + err.Error())
	}
	return t
}

var stage1Range = eventRange{
	Slug:      "cdl-major-1-tournament-2021",
	GameCode:  "CW",
	StartDate: mustUTC("2021-03-03 15:30"),
	EndDate:   mustUTC("2021-03-07 17:00"),
}

var slugs = map[string]uint{
	"cdl-major-1-tournament-2021": 42,
}

func TestFindTournamentForMatch_LateMatchOnEndDay(t *testing.T) {
	grandFinals := time.Date(2021, 3, 7, 23, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", grandFinals)
	assert.Equal(t, uint(42), got, "Grand Finals at 23:00 UTC on end-day should resolve to the correct tournament")
}

func TestFindTournamentForMatch_WinnersFinals(t *testing.T) {
	winnersFinals := time.Date(2021, 3, 7, 20, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", winnersFinals)
	assert.Equal(t, uint(42), got, "Winners Finals at 20:00 UTC on end-day should resolve to the correct tournament")
}

func TestFindTournamentForMatch_EarlyMatchWithinRange(t *testing.T) {
	midTournament := time.Date(2021, 3, 5, 18, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", midTournament)
	assert.Equal(t, uint(42), got)
}

func TestFindTournamentForMatch_MatchOnStartDay(t *testing.T) {
	startDay := time.Date(2021, 3, 3, 0, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", startDay)
	assert.Equal(t, uint(42), got, "Match on start day should be found")
}

func TestFindTournamentForMatch_MatchBeforeRange(t *testing.T) {
	before := time.Date(2021, 3, 2, 23, 59, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", before)
	assert.Equal(t, uint(0), got, "Match before tournament start should not match")
}

func TestFindTournamentForMatch_MatchAfterRange(t *testing.T) {
	after := time.Date(2021, 3, 8, 0, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", after)
	assert.Equal(t, uint(0), got, "Match after tournament end date should not match")
}

func TestFindTournamentForMatch_WrongGameCode(t *testing.T) {
	midTournament := time.Date(2021, 3, 5, 18, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "BO6", midTournament)
	assert.Equal(t, uint(0), got, "Wrong game code should not match")
}

func TestFindTournamentForMatch_ZeroTime(t *testing.T) {
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", time.Time{})
	assert.Equal(t, uint(0), got, "Zero time should return 0 (unparseable dates fall to fallback)")
}

func TestFindTournamentForMatch_MultipleRanges(t *testing.T) {
	stage2Range := eventRange{
		Slug:      "cdl-major-2-tournament-2021",
		GameCode:  "CW",
		StartDate: mustUTC("2021-04-07 14:00"),
		EndDate:   mustUTC("2021-04-11 17:00"),
	}
	allSlugs := map[string]uint{
		"cdl-major-1-tournament-2021": 42,
		"cdl-major-2-tournament-2021": 99,
	}

	ranges := []eventRange{stage1Range, stage2Range}
	s2GrandFinals := time.Date(2021, 4, 11, 22, 0, 0, 0, time.UTC)
	got := findTournamentForMatch(ranges, allSlugs, "CW", s2GrandFinals)
	assert.Equal(t, uint(99), got, "Late match on Stage 2 end-day should resolve to Stage 2")

	s1GrandFinals := time.Date(2021, 3, 7, 23, 0, 0, 0, time.UTC)
	got = findTournamentForMatch(ranges, allSlugs, "CW", s1GrandFinals)
	assert.Equal(t, uint(42), got, "Stage 1 Grand Finals should still resolve to Stage 1")
}
