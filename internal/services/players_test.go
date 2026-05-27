package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeEvent(name string, dates ...string) MatchEvent {
	matches := make([]MatchResult, len(dates))
	for i, d := range dates {
		matches[i] = MatchResult{Date: d}
	}
	return MatchEvent{Event: name, Matches: matches}
}

func TestSortEventsByDate_DateOverMatchID(t *testing.T) {
	events := []MatchEvent{
		makeEvent("CDL Major 1 2023", "0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"),
		makeEvent("Esports World Cup 2025", "2025-07-26T18:30:00Z", "2025-07-24T19:30:00Z"),
	}

	sortEventsByDate(events)

	assert.Equal(t, "Esports World Cup 2025", events[0].Event,
		"EWC 2025 must be first — real date beats zero date")
	assert.Equal(t, "CDL Major 1 2023", events[1].Event)
}

func TestSortEventsByDate_RealDatesOrdered(t *testing.T) {
	events := []MatchEvent{
		makeEvent("EWC 2024", "2024-08-17T00:00:00Z"),
		makeEvent("EWC 2025", "2025-07-26T18:30:00Z"),
	}

	sortEventsByDate(events)

	assert.Equal(t, "EWC 2025", events[0].Event)
	assert.Equal(t, "EWC 2024", events[1].Event)
}

func TestSortEventsByDate_EmptyEventsSink(t *testing.T) {
	events := []MatchEvent{
		makeEvent("Empty Event"),
		makeEvent("EWC 2025", "2025-07-26T18:30:00Z"),
	}

	sortEventsByDate(events)

	assert.Equal(t, "EWC 2025", events[0].Event)
	assert.Equal(t, "Empty Event", events[1].Event)
}

func TestSortEventsByDate_OrdersByRecentMatchFirst(t *testing.T) {
	events := []MatchEvent{
		makeEvent("B", "2024-08-01T00:00:00Z"),
		makeEvent("A", "2025-07-26T00:00:00Z"),
	}
	sortEventsByDate(events)
	assert.Equal(t, "A", events[0].Event)
	assert.Equal(t, "B", events[1].Event)
}

func TestSortEventsByDate_EmptyMatchesLast(t *testing.T) {
	events := []MatchEvent{
		makeEvent("empty"),
		makeEvent("has-match", "2025-01-01T00:00:00Z"),
	}
	sortEventsByDate(events)
	assert.Equal(t, "has-match", events[0].Event)
	assert.Equal(t, "empty", events[1].Event)
}

func TestSortEventsByDate_StableOnSingleEvent(t *testing.T) {
	events := []MatchEvent{
		makeEvent("only", "2025-07-26T00:00:00Z"),
	}
	sortEventsByDate(events)
	assert.Equal(t, "only", events[0].Event)
}

func TestCalculateKD(t *testing.T) {
	tests := []struct {
		name   string
		kills  int
		deaths int
		want   float64
	}{
		{"normal ratio", 10, 5, 2.0},
		{"equal kd", 5, 5, 1.0},
		{"below 1", 3, 6, 0.5},
		{"zero deaths returns 0", 10, 0, 0},
		{"both zero", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CalculateKD(tt.kills, tt.deaths))
		})
	}
}
