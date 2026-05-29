package main

// main.go — entry point for the CDL database seeder.
// This file only orchestrates the phases. All logic lives in the phase files.
//
// Seeder file structure:
//   config.go          — static config vars (which CSVs belong to which era)
//   types.go           — all CSV row structs and internal helper types
//   helpers.go         — shared utilities: date parsers, atoi/atof, DB resolution helpers
//   csv_readers.go     — one reader function per CSV file type
//   phase1_foundation.go — Franchises → CDL Teams → Non-CDL Teams → Players → Seasons → Tournaments
//   phase2_era_finals.go — Match + MatchMap + PlayerMapStats + PlayerMatchStats from era_finals/
//   phase3_enriched.go   — Same tables, EWC 2024/2025 + Major 1 2023 wiki data
//   phase4_season_stats.go — PlayerTournamentStats from *_player_stats.csv aggregates
//   phase5_transfers.go  — PlayerTransfer + unresolved_transfer_teams.csv report
//   phase6_bracket_patches.go — bracket_round + bracket_position backfill
//   phase7_rosters.go    — season-aware TeamRoster stints inferred from player_map_stats

import (
	"flag"
	"log"
	"os"

	"github.com/corbynfang/CDL-Website/internal/database"
)

func main() {
	reset := flag.Bool("reset", false, "truncate all seeder-owned tables before seeding (clean re-seed from empty)")
	autoYes := flag.Bool("yes", false, "skip the interactive confirmation for -reset (for non-interactive runs)")
	flag.Parse()

	database.ConnectDatabase()
	database.AutoMigrate()
	db := database.DB

	if *reset {
		if !confirmReset(*autoYes) {
			os.Exit(1)
		}
		if err := resetSeedTables(db); err != nil {
			log.Fatalf("reset failed: %v", err)
		}
	}

	log.Println("==> Cleanup: removing bad player records")
	cleanupBadPlayers(db)

	log.Println("==> Phase 1: Foundation (franchises, teams, players, seasons, tournaments)")
	franchiseMap := seedFranchises(db)
	teamLookup := seedCDLTeams(db, franchiseMap)
	mergeInto(teamLookup, seedNonCDLTeams(db))
	playerLookup := seedPlayers(db)
	seasonByCode := seedSeasons(db)
	tournamentBySlug, eventRanges := seedTournaments(db, seasonByCode)

	log.Println("==> Phase 2: era_finals match data (series, maps, player stats)")
	matchByBPID := seedEraFinals(db, teamLookup, playerLookup, seasonByCode, tournamentBySlug, eventRanges)
	_ = matchByBPID

	log.Println("==> Phase 3: Enriched match data (EWC 2024/2025, Major 1 2023 wiki)")
	seedEnrichedMatches(db, teamLookup, playerLookup, tournamentBySlug)

	log.Println("==> Phase 4: Season aggregate player stats")
	for _, cfg := range seasonStatConfigs {
		seasonID := seasonByCode[cfg.GameCode]
		if seasonID == 0 {
			continue
		}
		// Eras with a pre-aggregated CSV load from it; those without one (e.g. BO6)
		// derive their season totals from the per-map stats already seeded.
		if cfg.PlayerFile == "" {
			seedDerivedSeasonStats(db, cfg, seasonID)
			continue
		}
		seedSeasonStats(db, cfg, seasonID, teamLookup, playerLookup)
	}

	log.Println("==> Phase 5: Transfers (all 5 eras)")
	seedTransfers(db, teamLookup, playerLookup)

	log.Println("==> Phase 6: Bracket patches (bracket_round + bracket_position, all eras)")
	seedBracketPatches(db, teamLookup, tournamentBySlug)

	log.Println("==> Phase 7: Roster inference (season-aware stints from player_map_stats)")
	seedRosters(db)

	log.Println("==> Seeding complete.")
}
