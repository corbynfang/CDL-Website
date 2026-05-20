package main

// phase1_foundation.go — seeds the core reference data every other phase depends on.
// Order matters: Franchises → Teams → Players → Seasons → Tournaments.
// Each function returns a lookup map so downstream phases can resolve names to DB IDs
// without hitting the database on every row.

import (
	"log"
	"time"

	"github.com/corbynfang/CDL-Website/internal/database"
	"gorm.io/gorm"
)

// cleanupBadPlayers removes known CSV artifacts that aren't real CDL players.
func cleanupBadPlayers(db *gorm.DB) {
	for _, tag := range badGamertags {
		var p database.Player
		if err := db.Where("gamertag = ?", tag).First(&p).Error; err != nil {
			continue
		}
		db.Where("player_id = ?", p.ID).Delete(&database.PlayerMapStats{})
		db.Where("player_id = ?", p.ID).Delete(&database.PlayerMatchStats{})
		db.Where("player_id = ?", p.ID).Delete(&database.PlayerTournamentStats{})
		db.Delete(&p)
		log.Printf("Removed bad player: %s (id=%d)", tag, p.ID)
	}
}

// seedFranchises creates one Franchise row per unique CDL franchise slot (franchise_key).
// The franchise name is set to the most recent active branding so the display name is current.
// Returns franchise_key → franchise.ID.
func seedFranchises(db *gorm.DB) map[string]uint {
	rows := readBrandingCSV("database/cdl_team_branding_by_season.csv")

	type candidate struct {
		name      string
		validFrom string
		active    bool
	}
	best := map[string]candidate{}
	for _, r := range rows {
		if r.FranchiseKey == "" {
			continue
		}
		c, exists := best[r.FranchiseKey]
		isActive := r.ValidTo == ""
		if !exists || isActive || (!c.active && r.ValidFrom > c.validFrom) {
			best[r.FranchiseKey] = candidate{r.CanonicalTeamName, r.ValidFrom, isActive}
		}
	}

	franchiseMap := map[string]uint{}
	for key, c := range best {
		f := database.Franchise{FranchiseKey: key, Name: c.name, IsActive: c.active}
		db.Where("franchise_key = ?", key).FirstOrCreate(&f)
		if c.active {
			db.Model(&f).Update("name", c.name)
		}
		franchiseMap[key] = f.ID
	}
	log.Printf("Franchises seeded: %d", len(franchiseMap))
	return franchiseMap
}

// seedCDLTeams creates one Team row per row in cdl_team_branding_by_season.csv.
// Minnesota RØKKR and G2 Minnesota are two separate Team rows linked to the same Franchise.
// Returns a name→ID lookup covering both raw and canonical name variants.
func seedCDLTeams(db *gorm.DB, franchiseMap map[string]uint) map[string]uint {
	rows := readBrandingCSV("database/cdl_team_branding_by_season.csv")
	lookup := map[string]uint{}

	for _, r := range rows {
		franchiseID := franchiseMap[r.FranchiseKey]
		isActive := r.ValidTo == ""

		var validFrom, validTo *time.Time
		if t := parseFlexDate(r.ValidFrom); !t.IsZero() {
			validFrom = &t
		}
		if t := parseFlexDate(r.ValidTo); !t.IsZero() {
			validTo = &t
		}

		t := database.Team{
			Name:               r.CanonicalTeamName,
			Abbreviation:       makeAbbr(r.CanonicalTeamName),
			GameCode:           r.GameCode,
			IsCDLFranchise:     true,
			TeamClassification: "cdl_franchise",
			DoNotMerge:         false,
			ValidFrom:          validFrom,
			ValidTo:            validTo,
			IsActive:           isActive,
			Source:             "branding_csv",
		}
		if franchiseID != 0 {
			t.FranchiseID = &franchiseID
		}

		db.Where("name = ? AND source = ?", r.CanonicalTeamName, "branding_csv").FirstOrCreate(&t)
		lookup[r.CanonicalTeamName] = t.ID
		if r.RawTeamName != r.CanonicalTeamName && r.RawTeamName != "" {
			lookup[r.RawTeamName] = t.ID
		}
	}
	log.Printf("CDL teams seeded: %d lookup entries", len(lookup))
	return lookup
}

// seedNonCDLTeams creates Team rows for challenger orgs, academy teams, EWC international
// teams, and parent orgs from non_cdl_team_aliases_clean.csv.
// Returns a name→ID lookup merged into the CDL team lookup by main.go.
func seedNonCDLTeams(db *gorm.DB) map[string]uint {
	rows := readNonCDLCSV("database/non_cdl_team_aliases_clean.csv")
	lookup := map[string]uint{}

	for _, r := range rows {
		t := database.Team{
			Name:               r.CanonicalTeamName,
			Abbreviation:       makeAbbr(r.CanonicalTeamName),
			IsCDLFranchise:     false,
			TeamClassification: r.TeamClassification,
			DoNotMerge:         r.DoNotMerge,
			NeedsManualReview:  r.NeedsManualReview,
			Source:             "non_cdl_alias",
		}
		db.Where("name = ? AND source = ?", r.CanonicalTeamName, "non_cdl_alias").FirstOrCreate(&t)
		lookup[r.CanonicalTeamName] = t.ID
		if r.RawTeamName != r.CanonicalTeamName && r.RawTeamName != "" {
			lookup[r.RawTeamName] = t.ID
		}
	}
	log.Printf("Non-CDL teams seeded: %d lookup entries", len(lookup))
	return lookup
}

// seedPlayers reads player_aliases_clean.csv and creates one Player per unique canonical gamertag.
// Raw alias variants are also added to the lookup so transfer/stat rows that use alternate spellings still resolve.
// Returns gamertag → player.ID.
func seedPlayers(db *gorm.DB) map[string]uint {
	rows := readPlayerAliasCSV("database/player_aliases_clean.csv")
	lookup := map[string]uint{}

	for _, r := range rows {
		tag := r.CanonicalPlayerName
		if tag == "" {
			continue
		}
		isBad := false
		for _, bad := range badGamertags {
			if tag == bad {
				isBad = true
				break
			}
		}
		if isBad {
			continue
		}
		if existing, done := lookup[tag]; done {
			if r.PlayerName != tag && r.PlayerName != "" {
				lookup[r.PlayerName] = existing
			}
			continue
		}
		p := database.Player{Gamertag: tag}
		db.Where("gamertag = ?", tag).FirstOrCreate(&p)
		lookup[tag] = p.ID
		if r.PlayerName != tag && r.PlayerName != "" {
			lookup[r.PlayerName] = p.ID
		}
	}
	log.Printf("Players seeded: %d lookup entries", len(lookup))
	return lookup
}

// seedSeasons creates one Season row per CDL game era, keyed by game_code.
// Returns game_code → season.ID (e.g. "BO6" → 1).
func seedSeasons(db *gorm.DB) map[string]uint {
	type cfg struct {
		Code      string
		Name      string
		GameTitle string
		Year      int
	}
	cfgs := []cfg{
		{"BO6", "Black Ops 6 2024-25", "Black Ops 6", 2024},
		{"MW3", "Modern Warfare III 2023-24", "Modern Warfare III", 2023},
		{"MW2", "Modern Warfare II 2022-23", "Modern Warfare II", 2022},
		{"VG", "Vanguard 2021-22", "Vanguard", 2021},
		{"CW", "Black Ops Cold War 2020-21", "Black Ops Cold War", 2020},
	}
	byCode := map[string]uint{}
	for _, c := range cfgs {
		s := database.Season{
			Name:      c.Name,
			GameTitle: c.GameTitle,
			GameCode:  c.Code,
			StartDate: time.Date(c.Year, 9, 1, 0, 0, 0, 0, time.UTC),
			IsActive:  c.Code == "BO6",
		}
		db.Where("name = ?", c.Name).FirstOrCreate(&s)
		db.Model(&s).Update("game_code", c.Code)
		byCode[c.Code] = s.ID
	}
	log.Printf("Seasons seeded: %d", len(byCode))
	return byCode
}

// seedTournaments reads event_aliases_clean.csv and creates one Tournament per CDL event.
// Also builds the eventRanges slice used by findTournamentForMatch so match seeding can
// look up "which tournament does this date fall in?" without any extra DB queries.
// Returns (slug→tournament.ID, event ranges for date matching).
func seedTournaments(db *gorm.DB, seasonByCode map[string]uint) (map[string]uint, []eventRange) {
	rows := readEventAliasCSV("database/event_aliases_clean.csv")
	bySlug := map[string]uint{}
	var ranges []eventRange

	for _, r := range rows {
		if r.EventSlug == "" {
			continue
		}
		seasonID := seasonByCode[r.GameCode]
		if seasonID == 0 {
			continue
		}
		startDate := parseFlexDate(r.StartDate)
		var endDatePtr *time.Time
		if t := parseFlexDate(r.EndDate); !t.IsZero() {
			endDatePtr = &t
		}

		t := database.Tournament{
			SeasonID:         seasonID,
			Name:             r.CanonicalEventName,
			Slug:             r.EventSlug,
			TournamentType:   r.EventType,
			StartDate:        startDate,
			EndDate:          endDatePtr,
			SourceURL: r.SourceURL,
		}
		db.Where("slug = ? AND season_id = ?", r.EventSlug, seasonID).FirstOrCreate(&t)
		bySlug[r.EventSlug] = t.ID

		if !startDate.IsZero() {
			end := startDate.AddDate(0, 0, 30)
			if endDatePtr != nil {
				end = *endDatePtr
			}
			ranges = append(ranges, eventRange{
				Slug:      r.EventSlug,
				GameCode:  r.GameCode,
				StartDate: startDate,
				EndDate:   end,
			})
		}
	}
	log.Printf("Tournaments seeded: %d slugs", len(bySlug))
	return bySlug, ranges
}
