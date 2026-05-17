package main

// config.go — static configuration for every seeding phase.
// Each phase reads from these vars to know which files to process and what era they belong to.

// eraFinalsCfg describes the three CSV files that make up one era in database/era_finals/.
type eraFinalsCfg struct {
	GameCode   string
	SeriesFile string
	MapsFile   string
	StatsFile  string
}

// eraFinalsConfigs drives Phase 2. Order doesn't matter — each era is independent.
var eraFinalsConfigs = []eraFinalsCfg{
	{"BO6", "database/era_finals/bo6_series_final.csv", "database/era_finals/bo6_match_maps_final.csv", "database/era_finals/bo6_player_map_stats_final.csv"},
	{"CW", "database/era_finals/cw_series_final.csv", "database/era_finals/cw_match_maps_final.csv", "database/era_finals/cw_player_map_stats_final.csv"},
	{"MW2", "database/era_finals/mw2_series_final.csv", "database/era_finals/mw2_match_maps_final.csv", "database/era_finals/mw2_player_map_stats_final.csv"},
	{"MW3", "database/era_finals/mw3_series_final.csv", "database/era_finals/mw3_match_maps_final.csv", "database/era_finals/mw3_player_map_stats_final.csv"},
	{"VG", "database/era_finals/vg_series_final.csv", "database/era_finals/vg_match_maps_final.csv", "database/era_finals/vg_player_map_stats_final.csv"},
}

// seasonStatCfg drives Phase 4 — the season-level aggregate player stat CSVs.
type seasonStatCfg struct {
	GameCode   string
	PlayerFile string // empty string = no aggregate file for this era
	Name       string
	GameTitle  string
	StartYear  int
}

var seasonStatConfigs = []seasonStatCfg{
	{"BO6", "", "Black Ops 6 2024-25", "Black Ops 6", 2024},
	{"MW3", "database/cdl_mw3_player_stats.csv", "Modern Warfare III 2023-24", "Modern Warfare III", 2023},
	{"MW2", "database/cdl_mw2_players_stats.csv", "Modern Warfare II 2022-23", "Modern Warfare II", 2022},
	{"VG", "database/cdl_vanguard_players_stats.csv", "Vanguard 2021-22", "Vanguard", 2021},
	{"CW", "database/cdl_coldwar_players_stats.csv", "Black Ops Cold War 2020-21", "Black Ops Cold War", 2020},
}

// transferCfg drives Phase 5 — one entry per transfer CSV, tagged with its era.
type transferCfg struct {
	File     string
	GameCode string
	Season   string
}

var transferConfigs = []transferCfg{
	{"database/bo6_transfers.csv", "BO6", "Black Ops 6 2024-25"},
	{"database/cdl_mw3_transfers.csv", "MW3", "Modern Warfare III 2023-24"},
	{"database/cdl_mw2_transfers.csv", "MW2", "Modern Warfare II 2022-23"},
	{"database/cdl_vanguard_transfers.csv", "VG", "Vanguard 2021-22"},
	{"database/cdl_coldwar_transfers.csv", "CW", "Black Ops Cold War 2020-21"},
}

// badGamertags are CSV artifacts that are not real CDL players and should be purged.
var badGamertags = []string{"5aLDx"}
