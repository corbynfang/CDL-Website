package main

// match_map_played_test.go — regression test for the MatchMap.Played fix.
//
// MatchMap.Played used to carry gorm `default:true`, which made GORM omit a
// zero-value false on insert and let the DB default win — silently turning every
// DNP map into a played map (1,345 rows in the source CSVs). The default was
// removed; this test pins the behavior through the real seed insert path
// (CreateInBatches + OnConflict) so it can't regress.

import (
	"testing"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
)

func TestMatchMapPlayedFalseIsPersisted(t *testing.T) {
	db := rosterTx(t)

	mkSeason(t, db, 1, "BO6")
	mkTeam(t, db, 1, "OpTic Texas", "OTX", nil)
	mkTeam(t, db, 2, "Atlanta FaZe", "ATL", nil)
	mkTournament(t, db, 1, 1)
	mkMatch(t, db, 100, 1, 1, 2, time.Now())

	batch := []models.MatchMap{
		{MatchID: 100, MapNumber: 1, MapName: "Played Map", Played: true},
		{MatchID: 100, MapNumber: 2, MapName: "DNP Map", Played: false},
	}
	require.NoError(t, db.Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(batch, 500).Error)

	var got []models.MatchMap
	require.NoError(t, db.Where("match_id = ?", 100).Order("map_number").Find(&got).Error)
	require.Len(t, got, 2)
	require.True(t, got[0].Played, "played map must stay true")
	require.False(t, got[1].Played, "DNP map must stay false (no default:true corruption)")
}
