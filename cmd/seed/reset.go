package main

// reset.go — guarded "reset then seed" path.
//
// The seeder is additive (FirstOrCreate / OnConflict DoNothing), so re-running it
// over an existing database leaves stale rows in place — most importantly the old
// merged team rows and every fact already pointed at their IDs (match_maps,
// player_map_stats, player_match_stats, rosters, ...). Re-pointing those in place
// is fragile; the reliable fix is a clean re-seed from empty.
//
// resetSeedTables wipes every table the seeder owns and resets their identity
// sequences, so a subsequent seed run produces correct per-era team rows and links.
// It is destructive, so main only calls it behind the -reset flag and a confirmation
// (interactive y/N, or -yes for non-interactive runs).

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"gorm.io/gorm"
)

// resetTables lists every table populated by the seeder, ordered child → parent.
// TRUNCATE ... CASCADE makes the order non-critical, but listing children first
// keeps the intent (and any non-cascade dry runs) honest. This is the complete set
// the seeder owns; anything not here (e.g. data added outside the seeder) is left
// untouched, but note CASCADE will still clear rows in other tables that hold a
// foreign key into this set.
var resetTables = []string{
	"team_rosters",
	"player_map_stats",
	"player_match_stats",
	"player_tournament_stats",
	"team_tournament_stats",
	"player_transfers",
	"coaches",
	"match_maps",
	"matches",
	"tournaments",
	"teams",
	"players",
	"seasons",
	"franchises",
}

// resetSeedTables truncates every seeder-owned table in a single statement and
// restarts identity sequences so re-seeded rows get fresh, contiguous IDs.
func resetSeedTables(db *gorm.DB) error {
	stmt := fmt.Sprintf(
		"TRUNCATE TABLE %s RESTART IDENTITY CASCADE",
		strings.Join(resetTables, ", "),
	)
	log.Printf("==> Reset: truncating %d tables (RESTART IDENTITY CASCADE)", len(resetTables))
	if err := db.Exec(stmt).Error; err != nil {
		return fmt.Errorf("truncate failed: %w", err)
	}
	log.Println("==> Reset: tables cleared")
	return nil
}

// describeTarget returns a credential-free "host/dbname" label for the DATABASE_URL,
// so the confirmation prompt makes the blast radius obvious without leaking secrets.
func describeTarget() string {
	raw := os.Getenv("DATABASE_URL")
	if raw == "" {
		return "(DATABASE_URL unset)"
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "(unparseable DATABASE_URL)"
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		dbName = "(default db)"
	}
	return u.Host + "/" + dbName
}

// confirmReset gates the destructive truncate. With -yes it returns immediately;
// otherwise it prints the target and table list and requires a literal "yes".
func confirmReset(autoYes bool) bool {
	target := describeTarget()
	if autoYes {
		log.Printf("==> Reset: -yes supplied, proceeding against %s", target)
		return true
	}

	fmt.Printf("\nThis will TRUNCATE the following tables on %s:\n", target)
	for _, t := range resetTables {
		fmt.Printf("  - %s\n", t)
	}
	fmt.Print("\nAll seeded data will be deleted and re-created from CSV.\nType 'yes' to continue: ")

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(line)) != "yes" {
		fmt.Println("Aborted — no changes made.")
		return false
	}
	return true
}
