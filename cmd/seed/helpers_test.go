package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mustUTC parses "2006-01-02 15:04" as UTC. Panics on bad input — test helper only.
func mustUTC(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04", s, time.UTC)
	if err != nil {
		panic("mustUTC: " + err.Error())
	}
	return t
}

// stage1Range mirrors how event_aliases_clean.csv describes CDL Major 1 Tournament 2021:
// start 2021-03-03 3:30 pm, end 2021-03-07 5:00 pm (both UTC after parseFlexDate).
var stage1Range = eventRange{
	Slug:      "cdl-major-1-tournament-2021",
	GameCode:  "CW",
	StartDate: mustUTC("2021-03-03 15:30"),
	EndDate:   mustUTC("2021-03-07 17:00"),
}

var slugs = map[string]uint{
	"cdl-major-1-tournament-2021": 42,
}

// TestFindTournamentForMatch_LateMatchOnEndDay is the regression test for the root cause
// of the mapless-match bug: Grand Finals and other finals-day matches start after 5 PM
// UTC, which was past the tournament's CSV end-time. The seeder placed them in the
// fallback "unmatched" bucket instead of the correct tournament, causing Phase 6 to
// insert a duplicate mapless stub.
func TestFindTournamentForMatch_LateMatchOnEndDay(t *testing.T) {
	// Stage 1 Grand Finals: 2021-03-07T23:00 UTC — 6 hours after end-time of 17:00 UTC.
	grandFinals := time.Date(2021, 3, 7, 23, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", grandFinals)
	assert.Equal(t, uint(42), got, "Grand Finals at 23:00 UTC on end-day should resolve to the correct tournament")
}

func TestFindTournamentForMatch_WinnersFinals(t *testing.T) {
	// Winners Finals on the same day: 2021-03-07T20:00 UTC.
	winnersFinals := time.Date(2021, 3, 7, 20, 0, 0, 0, time.UTC)
	got := findTournamentForMatch([]eventRange{stage1Range}, slugs, "CW", winnersFinals)
	assert.Equal(t, uint(42), got, "Winners Finals at 20:00 UTC on end-day should resolve to the correct tournament")
}

func TestFindTournamentForMatch_EarlyMatchWithinRange(t *testing.T) {
	// A normal mid-tournament match well within the range.
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
	// A day after the end date — should not match.
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

// TestFindTournamentForMatch_AllRealEraFinals reads the actual CSV data that the seeder
// uses in production and asserts every match resolves to a real tournament (not the
// fallback 0). A failure here means a new CSV has dates outside all known event ranges.
func TestFindTournamentForMatch_AllRealEraFinals(t *testing.T) {
	aliasRows := readEventAliasCSV("../../database/event_aliases_clean.csv")
	if len(aliasRows) == 0 {
		t.Skip("event_aliases_clean.csv not available (gitignored proprietary data)")
	}

	// Build the same bySlug and ranges maps that the real seeder builds.
	bySlug := map[string]uint{}
	var ranges []eventRange
	for i, a := range aliasRows {
		slug := strings.TrimSpace(a.EventSlug)
		if slug == "" {
			continue
		}
		id := uint(i + 1) // synthetic IDs — we only care about non-zero
		bySlug[slug] = id
		startT := parseFlexDate(a.StartDate)
		endT := parseFlexDate(a.EndDate)
		if startT.IsZero() || endT.IsZero() {
			continue
		}
		ranges = append(ranges, eventRange{
			Slug:      slug,
			GameCode:  strings.TrimSpace(a.GameCode),
			StartDate: startT,
			EndDate:   endT,
		})
	}
	require.NotEmpty(t, ranges, "no valid event ranges parsed from event_aliases_clean.csv")

	// One sub-test per era so failures are attributed correctly.
	for _, era := range eraFinalsConfigs {
		era := era
		t.Run(era.GameCode, func(t *testing.T) {
			path := fmt.Sprintf("../../%s", era.SeriesFile)
			rows := readSeriesCSV(path)
			require.NotEmpty(t, rows, "no series rows in %s", path)

			var missed []string
			for _, s := range rows {
				matchTime := parseISOTime(s.MatchDatetime)
				if matchTime.IsZero() {
					// Unparseable date — skip; zero-time handling is covered by TestFindTournamentForMatch_ZeroTime.
					continue
				}
				id := findTournamentForMatch(ranges, bySlug, era.GameCode, matchTime)
				if id == 0 {
					missed = append(missed, fmt.Sprintf("match %d @ %s", s.MatchID, matchTime.UTC().Format(time.RFC3339)))
				}
			}
			assert.Empty(t, missed,
				"%s: %d match(es) fell outside all event ranges (would land in fallback bucket):\n  %s",
				era.GameCode, len(missed), strings.Join(missed, "\n  "))
		})
	}
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

	// Stage 2 Grand Finals late on its end day.
	s2GrandFinals := time.Date(2021, 4, 11, 22, 0, 0, 0, time.UTC)
	got := findTournamentForMatch(ranges, allSlugs, "CW", s2GrandFinals)
	assert.Equal(t, uint(99), got, "Late match on Stage 2 end-day should resolve to Stage 2")

	// Stage 1 Grand Finals should still resolve to Stage 1.
	s1GrandFinals := time.Date(2021, 3, 7, 23, 0, 0, 0, time.UTC)
	got = findTournamentForMatch(ranges, allSlugs, "CW", s1GrandFinals)
	assert.Equal(t, uint(42), got, "Stage 1 Grand Finals should still resolve to Stage 1")
}
