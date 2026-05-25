package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

func seedTransfers(db *gorm.DB, teamLookup map[string]uint, playerLookup map[string]uint) {
	report := map[string]*unresolvedTeamEntry{}
	total := 0

	for _, cfg := range transferConfigs {
		rows := readTransferCSV(cfg.File)
		for _, r := range rows {
			playerID := resolvePlayer(r.Player, playerLookup, db)
			if playerID == 0 {
				log.Printf("[transfers/%s] WARN: unknown player %q — skipping", cfg.GameCode, r.Player)
				continue
			}

			transferDate := parseTransferDate(r.Date)
			fromID := resolveTransferTeam(r.FromTeam, teamLookup, db, report, transferDate)
			toID := resolveTransferTeam(r.ToTeam, teamLookup, db, report, transferDate)

			var fromPtr, toPtr *uint
			if fromID != 0 {
				fromPtr = &fromID
			}
			if toID != 0 {
				toPtr = &toID
			}

			xfer := database.PlayerTransfer{
				PlayerID:        playerID,
				FromTeamID:      fromPtr,
				ToTeamID:        toPtr,
				TransferDate:    transferDate,
				TransferType:    r.TransferType,
				Role:            r.Role,
				GameCode:        cfg.GameCode,
				Season:          cfg.Season,
				RawFromTeamName: r.FromTeam,
				RawToTeamName:   r.ToTeam,
			}
			db.Where(
				"player_id = ? AND transfer_date = ? AND raw_from_team_name = ? AND raw_to_team_name = ?",
				playerID, transferDate, r.FromTeam, r.ToTeam,
			).FirstOrCreate(&xfer)
			total++
		}
	}

	log.Printf("Transfers seeded: %d rows", total)
	writeTransferReport(report)
}

func resolveTransferTeam(
	rawName string,
	teamLookup map[string]uint,
	db *gorm.DB,
	report map[string]*unresolvedTeamEntry,
	date time.Time,
) uint {
	rawName = strings.TrimSpace(rawName)
	if rawName == "" || rawName == "Free Agent" {
		return 0
	}

	if id, ok := teamLookup[rawName]; ok {
		recordResolution(report, rawName, date, "resolved_existing", false)
		return id
	}

	lower := strings.ToLower(rawName)
	for name, id := range teamLookup {
		if strings.ToLower(name) == lower {
			teamLookup[rawName] = id
			recordResolution(report, rawName, date, "resolved_existing", false)
			return id
		}
	}

	t := database.Team{
		Name:               rawName,
		Abbreviation:       makeAbbr(rawName),
		IsCDLFranchise:     false,
		TeamClassification: "unknown_challenger_or_regional",
		DoNotMerge:         true,
		NeedsManualReview:  true,
		Source:             "transfer_csv",
	}
	db.Where("name = ? AND source = ?", rawName, "transfer_csv").FirstOrCreate(&t)
	teamLookup[rawName] = t.ID
	recordResolution(report, rawName, date, "auto_created", true)
	return t.ID
}

func recordResolution(m map[string]*unresolvedTeamEntry, name string, date time.Time, status string, needsReview bool) {
	e, ok := m[name]
	if !ok {
		m[name] = &unresolvedTeamEntry{
			RawName:            name,
			Count:              1,
			FirstSeen:          date,
			LastSeen:           date,
			SuggestedCanonical: name,
			ResolutionStatus:   status,
			NeedsManualReview:  needsReview,
		}
		return
	}
	e.Count++
	if !date.IsZero() && (e.FirstSeen.IsZero() || date.Before(e.FirstSeen)) {
		e.FirstSeen = date
	}
	if !date.IsZero() && date.After(e.LastSeen) {
		e.LastSeen = date
	}
}

func writeTransferReport(report map[string]*unresolvedTeamEntry) {
	if len(report) == 0 {
		return
	}
	f, err := os.Create("database/unresolved_transfer_teams.csv")
	if err != nil {
		log.Printf("WARN: could not write transfer report: %v", err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	_ = w.Write([]string{
		"raw_team_name", "occurrence_count",
		"first_seen_transfer_date", "last_seen_transfer_date",
		"suggested_canonical_name", "resolution_status", "needs_manual_review",
	})
	for _, e := range report {
		_ = w.Write([]string{
			e.RawName,
			strconv.Itoa(e.Count),
			e.FirstSeen.Format("2006-01-02"),
			e.LastSeen.Format("2006-01-02"),
			e.SuggestedCanonical,
			e.ResolutionStatus,
			strconv.FormatBool(e.NeedsManualReview),
		})
	}
	w.Flush()
	log.Printf("Transfer team report written: database/unresolved_transfer_teams.csv (%d entries)", len(report))
}
