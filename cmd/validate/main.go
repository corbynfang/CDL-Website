package main

import (
	"fmt"
	"os"
)

// eraFiles groups the three CSV paths that belong to one CDL game era.
type eraFiles struct {
	gameCode   string
	seriesFile string
	mapsFile   string
	statsFile  string
}

var eras = []eraFiles{
	{"BO6", "database/era_finals/bo6_series_final.csv", "database/era_finals/bo6_match_maps_final.csv", "database/era_finals/bo6_player_map_stats_final.csv"},
	{"CW", "database/era_finals/cw_series_final.csv", "database/era_finals/cw_match_maps_final.csv", "database/era_finals/cw_player_map_stats_final.csv"},
	{"MW2", "database/era_finals/mw2_series_final.csv", "database/era_finals/mw2_match_maps_final.csv", "database/era_finals/mw2_player_map_stats_final.csv"},
	{"MW3", "database/era_finals/mw3_series_final.csv", "database/era_finals/mw3_match_maps_final.csv", "database/era_finals/mw3_player_map_stats_final.csv"},
	{"VG", "database/era_finals/vg_series_final.csv", "database/era_finals/vg_match_maps_final.csv", "database/era_finals/vg_player_map_stats_final.csv"},
}

var transferFiles = []string{
	"database/bo6_transfers.csv",
	"database/cdl_mw3_transfers.csv",
	"database/cdl_mw2_transfers.csv",
	"database/cdl_vanguard_transfers.csv",
	"database/cdl_coldwar_transfers.csv",
}

func main() {
	var issues []Issue
	filesChecked := 0

	// Phase 1 — foundation CSVs
	filesChecked += 4
	issues = append(issues, validateNonCDLTeams("database/non_cdl_team_aliases_clean.csv")...)
	issues = append(issues, validatePlayerAliases("database/player_aliases_clean.csv")...)
	issues = append(issues, validateEventAliases("database/event_aliases_clean.csv")...)
	issues = append(issues, validateBranding("database/cdl_team_branding_by_season.csv")...)

	// Phase 2 — era finals (15 files: series + maps + stats × 5 eras)
	for _, era := range eras {
		filesChecked += 3
		issues = append(issues, validateSeriesCSV(era.seriesFile)...)
		issues = append(issues, validateMapsCSV(era.mapsFile)...)
		issues = append(issues, validatePlayerStatsCSV(era.statsFile)...)
		issues = append(issues, crossReferenceEra(era)...)
	}

	// Phase 5 — transfer CSVs
	for _, f := range transferFiles {
		filesChecked++
		issues = append(issues, validateTransferCSV(f)...)
	}

	// Print every issue with its file and line number
	errors, warns := 0, 0
	for _, iss := range issues {
		fmt.Println(iss)
		if iss.level == levelError {
			errors++
		} else {
			warns++
		}
	}

	fmt.Printf("\n%d files checked — %d error(s)  %d warning(s)\n", filesChecked, errors, warns)
	if errors > 0 {
		fmt.Println("Fix errors before running the seeder.")
		os.Exit(1)
	}
	fmt.Println("All checks passed. Safe to seed.")
}
