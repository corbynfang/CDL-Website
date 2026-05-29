package main

// phase4_derived.go — derives season-level aggregate player stats from per-map data.
//
// Most eras ship a pre-aggregated *_player_stats.csv (see phase4_season_stats.go).
// BO6 has no such file (seasonStatCfg.PlayerFile == ""), so without this step its
// player_tournament_stats table stays empty and every view backed by it — the K/D
// rankings leaderboard in particular — returns nothing.
//
// We already seed the full per-map truth (player_map_stats) for BO6, so we roll that
// up into the same season-summary representation the CSV eras use: one row per player,
// attached to a virtual "<GAME>-season-stats" tournament. Only the overall totals
// (kills/deaths/assists/damage + K/D, KDA, maps) are derivable here; the HP/SND/CTL
// mode splits have no per-map source and are left at zero.
//
// The aggregation runs as one server-side INSERT ... SELECT so the same statement can
// be used verbatim as a targeted prod backfill, and it is idempotent
// (ON CONFLICT DO NOTHING) so re-running a full seed is safe.

import (
	"log"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

// derivedSeasonStatsSQL aggregates player_map_stats for one season into
// player_tournament_stats rows pointed at the given summary tournament.
// Params (in order): seasonID, tournamentID.
//   - dom_team: each player's team for the season = the team they played the most
//     maps for (ties broken by lower team_id for determinism).
//   - totals:   season-wide kills/deaths/assists/damage and map count per player.
// Players with no kills and no deaths are skipped to keep the leaderboard clean.
const derivedSeasonStatsSQL = `
WITH season_maps AS (
	SELECT pms.player_id, pms.team_id, pms.kills, pms.deaths, pms.assists, pms.damage
	FROM player_map_stats pms
	JOIN matches m      ON m.id = pms.match_id
	JOIN tournaments tn ON tn.id = m.tournament_id
	WHERE tn.season_id = ?
),
dom_team AS (
	SELECT DISTINCT ON (player_id) player_id, team_id
	FROM (
		SELECT player_id, team_id, COUNT(*) AS maps
		FROM season_maps
		GROUP BY player_id, team_id
	) per_team
	ORDER BY player_id, maps DESC, team_id
),
totals AS (
	SELECT player_id,
		SUM(kills)   AS k,
		SUM(deaths)  AS d,
		SUM(assists) AS a,
		SUM(damage)  AS dmg,
		COUNT(*)     AS maps
	FROM season_maps
	GROUP BY player_id
	HAVING SUM(kills) > 0 OR SUM(deaths) > 0
)
INSERT INTO player_tournament_stats
	(player_id, team_id, tournament_id, total_kills, total_deaths, total_assists,
	 total_damage, kd_ratio, kda_ratio, overall_maps)
SELECT t.player_id, dt.team_id, ?, t.k, t.d, t.a, t.dmg,
	CASE WHEN t.d > 0 THEN ROUND(t.k::decimal / t.d, 3) ELSE 0 END,
	CASE WHEN t.d > 0 THEN ROUND((t.k + t.a)::decimal / t.d, 3) ELSE 0 END,
	t.maps
FROM totals t
JOIN dom_team dt ON dt.player_id = t.player_id
ON CONFLICT (player_id, tournament_id) DO NOTHING`

// seedDerivedSeasonStats populates player_tournament_stats for an era that has no
// pre-aggregated CSV, by rolling up its per-map stats. It mirrors the season-summary
// tournament that seedSeasonStats creates so all eras share one representation.
func seedDerivedSeasonStats(db *gorm.DB, cfg seasonStatCfg, seasonID uint) {
	summaryTournament := models.Tournament{
		SeasonID:       seasonID,
		Name:           cfg.Name + " — Season Stats",
		Slug:           cfg.GameCode + "-season-stats",
		TournamentType: "season_summary",
		StartDate:      time.Date(cfg.StartYear, 6, 1, 0, 0, 0, 0, time.UTC),
	}
	db.Where("slug = ? AND season_id = ?", summaryTournament.Slug, seasonID).FirstOrCreate(&summaryTournament)

	res := db.Exec(derivedSeasonStatsSQL, seasonID, summaryTournament.ID)
	if res.Error != nil {
		log.Printf("[%s] derived season stats failed: %v", cfg.GameCode, res.Error)
		return
	}
	log.Printf("[%s] season stats derived from map data: %d rows", cfg.GameCode, res.RowsAffected)
}
