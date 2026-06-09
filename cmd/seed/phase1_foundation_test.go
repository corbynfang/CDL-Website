package main

import (
	"testing"

	"github.com/corbynfang/CDL-Website/internal/models"
	"github.com/stretchr/testify/require"
)

func TestResolveTeamID_GameAwareWithBareFallback(t *testing.T) {
	lookup := map[string]uint{
		teamKey("London Royal Ravens", "CW"):  1,
		teamKey("London Royal Ravens", "VG"):  2,
		teamKey("London Royal Ravens", "MW2"): 3,
		"London Royal Ravens":                 3, // bare fallback = last era
		teamKey("Toronto Ultra", "BO6"):       9,
		"Toronto Ultra":                       9,
	}

	require.Equal(t, uint(1), resolveTeamID(lookup, "London Royal Ravens", "CW"))
	require.Equal(t, uint(2), resolveTeamID(lookup, "London Royal Ravens", "VG"))
	require.Equal(t, uint(3), resolveTeamID(lookup, "London Royal Ravens", "MW2"))
	require.Equal(t, uint(3), resolveTeamID(lookup, "London Royal Ravens", "BO6"))
	require.Equal(t, uint(9), resolveTeamID(lookup, "  Toronto Ultra  ", "BO6"))
	require.Equal(t, uint(0), resolveTeamID(lookup, "Nonexistent", "BO6"))
	require.Equal(t, uint(0), resolveTeamID(lookup, "", "BO6"))
}

func TestSeedCDLTeamRows_SplitsRoyalRavensPerGameEra(t *testing.T) {
	db := rosterTx(t)

	franchiseID := uint(1)
	mkFranchise(t, db, franchiseID, "royal-ravens")

	rows := []brandingRow{
		{GameCode: "CW", RawTeamName: "London Royal Ravens", CanonicalTeamName: "London Royal Ravens", FranchiseKey: "royal-ravens", ValidFrom: "2021-02-04", ValidTo: "2021-08-22"},
		{GameCode: "VG", RawTeamName: "London Royal Ravens", CanonicalTeamName: "London Royal Ravens", FranchiseKey: "royal-ravens", ValidFrom: "2022-02-04", ValidTo: "2022-07-31"},
		{GameCode: "MW2", RawTeamName: "London Royal Ravens", CanonicalTeamName: "London Royal Ravens", FranchiseKey: "royal-ravens", ValidFrom: "2022-12-15", ValidTo: "2023-08-20"},
		{GameCode: "MW3", RawTeamName: "Carolina Royal Ravens", CanonicalTeamName: "Carolina Royal Ravens", FranchiseKey: "royal-ravens", ValidFrom: "2024-01-12", ValidTo: "2024-08-25"},
		{GameCode: "BO6", RawTeamName: "Carolina Royal Ravens", CanonicalTeamName: "Carolina Royal Ravens", FranchiseKey: "royal-ravens", ValidFrom: "2024-12-06", ValidTo: "2025-08-31"},
	}

	lookup := seedCDLTeamRows(db, rows, map[string]uint{"royal-ravens": franchiseID})

	var teams []models.Team
	require.NoError(t, db.Where("franchise_id = ?", franchiseID).Find(&teams).Error)
	require.Len(t, teams, 5, "franchise should expose 5 era rows (LDN CW/VG/MW2 + CAR MW3/BO6)")

	byGame := map[string]models.Team{}
	for _, tm := range teams {
		byGame[tm.GameCode] = tm
	}
	for _, code := range []string{"CW", "VG", "MW2", "MW3", "BO6"} {
		require.Contains(t, byGame, code, "missing era row for %s", code)
	}
	require.Equal(t, "London Royal Ravens", byGame["CW"].Name)
	require.Equal(t, "London Royal Ravens", byGame["VG"].Name)
	require.Equal(t, "Carolina Royal Ravens", byGame["MW3"].Name)
	require.Equal(t, "Carolina Royal Ravens", byGame["BO6"].Name)
	require.NotEqual(t,
		resolveTeamID(lookup, "London Royal Ravens", "CW"),
		resolveTeamID(lookup, "London Royal Ravens", "VG"),
		"CW and VG must be separate London Royal Ravens rows")
	require.NotEqual(t,
		resolveTeamID(lookup, "Carolina Royal Ravens", "MW3"),
		resolveTeamID(lookup, "Carolina Royal Ravens", "BO6"),
		"MW3 and BO6 must be separate Carolina Royal Ravens rows")
	seedCDLTeamRows(db, rows, map[string]uint{"royal-ravens": franchiseID})
	var count int64
	require.NoError(t, db.Model(&models.Team{}).Where("franchise_id = ?", franchiseID).Count(&count).Error)
	require.EqualValues(t, 5, count, "second seed run must not duplicate era rows")
}
