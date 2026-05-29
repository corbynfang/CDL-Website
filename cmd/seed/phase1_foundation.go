package main

import (
	"log"
	"time"

	"github.com/corbynfang/CDL-Website/internal/models"
	"gorm.io/gorm"
)

func cleanupBadPlayers(db *gorm.DB) {
	for _, tag := range badGamertags {
		var p models.Player
		if err := db.Where("gamertag = ?", tag).First(&p).Error; err != nil {
			continue
		}
		db.Where("player_id = ?", p.ID).Delete(&models.PlayerMapStats{})
		db.Where("player_id = ?", p.ID).Delete(&models.PlayerMatchStats{})
		db.Where("player_id = ?", p.ID).Delete(&models.PlayerTournamentStats{})
		db.Delete(&p)
		log.Printf("Removed bad player: %s (id=%d)", tag, p.ID)
	}
}

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
		f := models.Franchise{FranchiseKey: key, Name: c.name, IsActive: c.active}
		db.Where("franchise_key = ?", key).FirstOrCreate(&f)
		if c.active {
			db.Model(&f).Update("name", c.name)
		}
		franchiseMap[key] = f.ID
	}
	log.Printf("Franchises seeded: %d", len(franchiseMap))
	return franchiseMap
}

func seedCDLTeams(db *gorm.DB, franchiseMap map[string]uint) map[string]uint {
	rows := readBrandingCSV("database/cdl_team_branding_by_season.csv")
	return seedCDLTeamRows(db, rows, franchiseMap)
}

// seedCDLTeamRows is the testable core of seedCDLTeams. Each branding row is one
// (team, game) era, and we create one team row per era — keyed on
// (name, game_code, source). Keying on name alone (the previous behavior) merged
// every game-era of a franchise name into a single row: London Royal Ravens'
// CW/VG/MW2 rows collapsed into one, and Carolina Royal Ravens' MW3/BO6 into one,
// which both broke the era list and leaked the latest season's roster onto older
// eras. The returned lookup carries a game-aware key per era plus a bare-name key
// (last era wins) for callers without game context.
func seedCDLTeamRows(db *gorm.DB, rows []brandingRow, franchiseMap map[string]uint) map[string]uint {
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

		t := models.Team{
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

		db.Where("name = ? AND game_code = ? AND source = ?", r.CanonicalTeamName, r.GameCode, "branding_csv").FirstOrCreate(&t)

		lookup[teamKey(r.CanonicalTeamName, r.GameCode)] = t.ID
		lookup[r.CanonicalTeamName] = t.ID // bare fallback (last era wins)
		if r.RawTeamName != r.CanonicalTeamName && r.RawTeamName != "" {
			lookup[teamKey(r.RawTeamName, r.GameCode)] = t.ID
			lookup[r.RawTeamName] = t.ID
		}
	}
	log.Printf("CDL teams seeded: %d lookup entries", len(lookup))
	return lookup
}

func seedNonCDLTeams(db *gorm.DB) map[string]uint {
	rows := readNonCDLCSV("database/non_cdl_team_aliases_clean.csv")
	lookup := map[string]uint{}

	for _, r := range rows {
		t := models.Team{
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
		p := models.Player{Gamertag: tag}
		db.Where("gamertag = ?", tag).FirstOrCreate(&p)
		lookup[tag] = p.ID
		if r.PlayerName != tag && r.PlayerName != "" {
			lookup[r.PlayerName] = p.ID
		}
	}
	log.Printf("Players seeded: %d lookup entries", len(lookup))
	return lookup
}

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
		s := models.Season{
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

		t := models.Tournament{
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
