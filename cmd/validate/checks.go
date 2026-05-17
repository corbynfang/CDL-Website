package main

// checks.go — one validation function per CSV file type.
// Every function returns []Issue so the caller collects ALL problems across
// all files before printing, rather than stopping at the first error.
//
// Two issue levels:
//   ERROR — seeding will break or produce wrong data; must fix before seeding.
//   WARN  — suspicious but won't crash the seeder; review before deploying.

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ─── Issue type ───────────────────────────────────────────────────────────────

type issueLevel int

const (
	levelError issueLevel = iota
	levelWarn
)

// Issue is one problem found in a CSV file, with the file path and line number
// so you can jump directly to the problem.
type Issue struct {
	level issueLevel
	file  string
	line  int // 1-based; 0 = file-level problem (e.g. missing file)
	msg   string
}

func (i Issue) String() string {
	tag := "ERROR"
	if i.level == levelWarn {
		tag = "WARN "
	}
	if i.line > 0 {
		return fmt.Sprintf("[%s] %s  line %d — %s", tag, i.file, i.line, i.msg)
	}
	return fmt.Sprintf("[%s] %s — %s", tag, i.file, i.msg)
}

func errorf(file string, line int, format string, args ...any) Issue {
	return Issue{levelError, file, line, fmt.Sprintf(format, args...)}
}

func warnf(file string, line int, format string, args ...any) Issue {
	return Issue{levelWarn, file, line, fmt.Sprintf(format, args...)}
}

// ─── Raw CSV loader ───────────────────────────────────────────────────────────

// loadRaw opens a CSV file and returns (headers, data rows, issues).
// FieldsPerRecord=-1 so it never stops on field count mismatches — instead,
// rows with the wrong count are reported as errors and excluded from the returned
// rows slice so downstream checks don't index out of bounds.
func loadRaw(path string) ([]string, [][]string, []Issue) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, []Issue{errorf(path, 0, "cannot open file: %v", err)}
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, []Issue{errorf(path, 0, "CSV parse error: %v", err)}
	}
	if len(records) < 2 {
		return nil, nil, []Issue{errorf(path, 0, "file is empty or header-only")}
	}

	headers := records[0]
	expected := len(headers)
	var issues []Issue
	var rows [][]string

	for i, rec := range records[1:] {
		line := i + 2 // +1 for 0-index, +1 because header is line 1
		if len(rec) != expected {
			issues = append(issues, errorf(path, line,
				"wrong field count: expected %d got %d — unquoted comma in a text field?",
				expected, len(rec)))
			continue // skip malformed row so downstream checks don't panic
		}
		rows = append(rows, rec)
	}
	return headers, rows, issues
}

// ─── Header and column helpers ────────────────────────────────────────────────

// requireHeaders checks that every name in required appears in the header row.
func requireHeaders(path string, headers []string, required []string) []Issue {
	set := map[string]bool{}
	for _, h := range headers {
		set[strings.ToLower(strings.TrimSpace(h))] = true
	}
	var issues []Issue
	for _, req := range required {
		if !set[strings.ToLower(req)] {
			issues = append(issues, errorf(path, 1, "missing required column %q", req))
		}
	}
	return issues
}

// colIdx returns the 0-based index of a column name (case-insensitive), or -1.
func colIdx(headers []string, name string) int {
	name = strings.ToLower(strings.TrimSpace(name))
	for i, h := range headers {
		if strings.ToLower(strings.TrimSpace(h)) == name {
			return i
		}
	}
	return -1
}

// cell returns the trimmed value at a column index, or "" if out of range.
func cell(row []string, i int) string {
	if i < 0 || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

// ─── Date parsers ─────────────────────────────────────────────────────────────

func isValidISO(s string) bool {
	_, e1 := time.Parse(time.RFC3339, s)
	_, e2 := time.Parse("2006-01-02T15:04:05", s)
	return e1 == nil || e2 == nil
}

func isValidFlexDate(s string) bool {
	for _, f := range []string{"2006-01-02 3:04 pm", "2006-01-02 15:04", "2006-01-02"} {
		if _, err := time.Parse(f, strings.ToLower(s)); err == nil {
			return true
		}
	}
	return false
}

func isValidTransferDate(s string) bool {
	for _, f := range []string{"Jan 2 2006", "Jan 02 2006"} {
		if _, err := time.Parse(f, s); err == nil {
			return true
		}
	}
	return false
}

// ─── Validators ───────────────────────────────────────────────────────────────

// validateNonCDLTeams validates database/non_cdl_team_aliases_clean.csv.
// This is the file that broke previously (Project 7, line 19 — unquoted comma in notes).
// The loadRaw field-count check catches that class of bug automatically.
func validateNonCDLTeams(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"raw_team_name", "canonical_team_name", "team_slug",
		"is_cdl_franchise", "team_classification", "do_not_merge",
		"needs_manual_review", "notes",
	})...)

	canonicalCol := colIdx(headers, "canonical_team_name")
	seen := map[string]int{}

	for i, row := range rows {
		line := i + 2
		name := cell(row, canonicalCol)
		if name == "" {
			issues = append(issues, errorf(path, line, "canonical_team_name is empty"))
			continue
		}
		if prev, ok := seen[name]; ok {
			issues = append(issues, warnf(path, line, "duplicate canonical_team_name %q (also on line %d)", name, prev))
		} else {
			seen[name] = line
		}
	}
	return issues
}

// validatePlayerAliases validates database/player_aliases_clean.csv.
func validatePlayerAliases(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"player_name", "canonical_player_name",
	})...)

	canonicalCol := colIdx(headers, "canonical_player_name")
	for i, row := range rows {
		if cell(row, canonicalCol) == "" {
			issues = append(issues, errorf(path, i+2, "canonical_player_name is empty"))
		}
	}
	return issues
}

// validateEventAliases validates database/event_aliases_clean.csv.
func validateEventAliases(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"event_slug", "game_code", "canonical_event_name", "start_date",
	})...)

	slugCol := colIdx(headers, "event_slug")
	dateCol := colIdx(headers, "start_date")
	seen := map[string]int{}

	for i, row := range rows {
		line := i + 2
		slug := cell(row, slugCol)
		if slug == "" {
			issues = append(issues, errorf(path, line, "event_slug is empty"))
			continue
		}
		if prev, ok := seen[slug]; ok {
			issues = append(issues, errorf(path, line, "duplicate event_slug %q (also on line %d)", slug, prev))
		} else {
			seen[slug] = line
		}
		if d := cell(row, dateCol); d != "" && !isValidFlexDate(d) {
			issues = append(issues, warnf(path, line, "unparseable start_date %q", d))
		}
	}
	return issues
}

// validateBranding validates database/cdl_team_branding_by_season.csv.
func validateBranding(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"game_code", "raw_team_name", "canonical_team_name", "franchise_key",
	})...)

	gameCodeCol := colIdx(headers, "game_code")
	canonicalCol := colIdx(headers, "canonical_team_name")
	franchiseCol := colIdx(headers, "franchise_key")
	validCodes := map[string]bool{"BO6": true, "CW": true, "MW2": true, "MW3": true, "VG": true}

	for i, row := range rows {
		line := i + 2
		if cell(row, canonicalCol) == "" {
			issues = append(issues, errorf(path, line, "canonical_team_name is empty"))
		}
		if cell(row, franchiseCol) == "" {
			issues = append(issues, warnf(path, line, "franchise_key is empty for team %q", cell(row, canonicalCol)))
		}
		if gc := cell(row, gameCodeCol); gc != "" && !validCodes[gc] {
			issues = append(issues, warnf(path, line, "unrecognised game_code %q", gc))
		}
	}
	return issues
}

// validateSeriesCSV validates an era_finals *_series_final.csv file.
func validateSeriesCSV(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"match_id", "match_datetime", "team_a_name", "team_b_name",
		"team_a_score", "team_b_score", "winner_name", "series_format",
	})...)

	matchIDCol := colIdx(headers, "match_id")
	datetimeCol := colIdx(headers, "match_datetime")
	teamACol := colIdx(headers, "team_a_name")
	teamBCol := colIdx(headers, "team_b_name")
	formatCol := colIdx(headers, "series_format")
	validFormats := map[string]bool{"BO3": true, "BO5": true, "BO7": true, "BO9": true, "": true}
	seen := map[string]int{}

	for i, row := range rows {
		line := i + 2
		id := cell(row, matchIDCol)
		if id == "" || id == "0" {
			issues = append(issues, errorf(path, line, "match_id is missing or zero"))
		} else if prev, ok := seen[id]; ok {
			issues = append(issues, errorf(path, line, "duplicate match_id %s (also on line %d)", id, prev))
		} else {
			seen[id] = line
		}
		if d := cell(row, datetimeCol); d != "" && !isValidISO(d) {
			issues = append(issues, warnf(path, line, "unparseable match_datetime %q", d))
		}
		if cell(row, teamACol) == "" {
			issues = append(issues, errorf(path, line, "team_a_name is empty"))
		}
		if cell(row, teamBCol) == "" {
			issues = append(issues, errorf(path, line, "team_b_name is empty"))
		}
		if f := cell(row, formatCol); !validFormats[f] {
			issues = append(issues, warnf(path, line, "unexpected series_format %q (expected BO3/BO5/BO7/BO9)", f))
		}
	}
	return issues
}

// validateMapsCSV validates an era_finals *_match_maps_final.csv file.
func validateMapsCSV(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"match_id", "map_number", "map_name", "mode_name",
	})...)

	matchIDCol := colIdx(headers, "match_id")
	mapNumCol := colIdx(headers, "map_number")
	modeCol := colIdx(headers, "mode_name")
	seen := map[string]int{}
	validModes := map[string]bool{
		"Hardpoint": true, "Search and Destroy": true,
		"Search & Destroy": true, "Control": true,
	}

	for i, row := range rows {
		line := i + 2
		key := cell(row, matchIDCol) + ":" + cell(row, mapNumCol)
		if prev, ok := seen[key]; ok {
			issues = append(issues, errorf(path, line,
				"duplicate (match_id=%s, map_number=%s) — also on line %d",
				cell(row, matchIDCol), cell(row, mapNumCol), prev))
		} else {
			seen[key] = line
		}
		if n, err := strconv.Atoi(cell(row, mapNumCol)); err != nil || n < 1 || n > 9 {
			issues = append(issues, warnf(path, line, "map_number %q is outside expected range 1–9", cell(row, mapNumCol)))
		}
		if m := cell(row, modeCol); m != "" && !validModes[m] {
			issues = append(issues, warnf(path, line, "unrecognised mode_name %q", m))
		}
	}
	return issues
}

// validatePlayerStatsCSV validates an era_finals *_player_map_stats_final.csv file.
func validatePlayerStatsCSV(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"match_id", "map_number", "player_id", "player_tag", "kills", "deaths",
	})...)

	matchIDCol := colIdx(headers, "match_id")
	mapNumCol := colIdx(headers, "map_number")
	playerIDCol := colIdx(headers, "player_id")
	playerTagCol := colIdx(headers, "player_tag")
	killsCol := colIdx(headers, "kills")
	deathsCol := colIdx(headers, "deaths")
	seen := map[string]int{}

	for i, row := range rows {
		line := i + 2
		key := cell(row, matchIDCol) + ":" + cell(row, mapNumCol) + ":" + cell(row, playerIDCol)
		if prev, ok := seen[key]; ok {
			issues = append(issues, errorf(path, line,
				"duplicate (match_id, map_number, player_id) — also on line %d", prev))
		} else {
			seen[key] = line
		}
		if cell(row, playerTagCol) == "" {
			issues = append(issues, errorf(path, line, "player_tag is empty"))
		}
		if k, err := strconv.Atoi(cell(row, killsCol)); err == nil && k < 0 {
			issues = append(issues, warnf(path, line, "negative kills value: %d", k))
		}
		if d, err := strconv.Atoi(cell(row, deathsCol)); err == nil && d < 0 {
			issues = append(issues, warnf(path, line, "negative deaths value: %d", d))
		}
	}
	return issues
}

// validateTransferCSV validates a transfer CSV (any era).
func validateTransferCSV(path string) []Issue {
	headers, rows, issues := loadRaw(path)
	if headers == nil {
		return issues
	}
	issues = append(issues, requireHeaders(path, headers, []string{
		"date", "player", "from_team", "to_team", "transfer_type",
	})...)

	dateCol := colIdx(headers, "date")
	playerCol := colIdx(headers, "player")
	typeCol := colIdx(headers, "transfer_type")
	validTypes := map[string]bool{
		"Signing": true, "Transfer": true, "Release": true,
		"Retirement": true, "Role Change": true,
	}

	for i, row := range rows {
		line := i + 2
		if d := cell(row, dateCol); d != "" && !isValidTransferDate(d) {
			issues = append(issues, warnf(path, line, "unparseable date %q (expected: Jan 2 2006)", d))
		}
		if cell(row, playerCol) == "" {
			issues = append(issues, errorf(path, line, "player is empty"))
		}
		if t := cell(row, typeCol); t != "" && !validTypes[t] {
			issues = append(issues, warnf(path, line, "unrecognised transfer_type %q", t))
		}
	}
	return issues
}

// crossReferenceEra checks that every match_id in the maps and stats files
// appears in the series file. Orphaned rows point to matches that don't exist
// and will be silently skipped by the seeder — catch them here instead.
func crossReferenceEra(era eraFiles) []Issue {
	sHeaders, sRows, issues := loadRaw(era.seriesFile)
	if sHeaders == nil {
		return issues
	}

	// Build the set of known match IDs from the series file.
	matchIDCol := colIdx(sHeaders, "match_id")
	known := map[string]bool{}
	for _, row := range sRows {
		known[cell(row, matchIDCol)] = true
	}

	// checkOrphans reports the first occurrence of each unknown match_id in path.
	checkOrphans := func(path, idColName string) {
		h, rows, _ := loadRaw(path) // field-count issues already reported by validate*
		if h == nil {
			return
		}
		col := colIdx(h, idColName)
		reported := map[string]bool{}
		for i, row := range rows {
			id := cell(row, col)
			if !known[id] && !reported[id] {
				issues = append(issues, errorf(path, i+2,
					"match_id %s not found in %s", id, era.seriesFile))
				reported[id] = true
			}
		}
	}

	checkOrphans(era.mapsFile, "match_id")
	checkOrphans(era.statsFile, "match_id")
	return issues
}
