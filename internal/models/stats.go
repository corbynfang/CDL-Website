package models

import "time"

type PlayerMapStats struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	MatchID   uint `json:"match_id" gorm:"not null;uniqueIndex:idx_player_map_stat_unique"`
	MapNumber int  `json:"map_number" gorm:"not null;uniqueIndex:idx_player_map_stat_unique"`
	PlayerID  uint `json:"player_id" gorm:"not null;uniqueIndex:idx_player_map_stat_unique;index"`
	TeamID    uint `json:"team_id" gorm:"not null;index"`
	Kills   int     `json:"kills" gorm:"default:0"`
	Deaths  int     `json:"deaths" gorm:"default:0"`
	KDRatio float64 `json:"kd_ratio" gorm:"type:decimal(6,3);default:0"`
	Damage  int     `json:"damage" gorm:"default:0"`
	Assists int     `json:"assists" gorm:"default:0"`

	BPRating float64 `json:"bp_rating" gorm:"type:decimal(10,6);default:0"` // Only used for BreakingPoint Stats if found BPRating within database.

	HillTime             int `json:"hill_time" gorm:"default:0"`
	SndRounds            int `json:"snd_rounds" gorm:"default:0"`
	PlantCount           int `json:"plant_count" gorm:"default:0"`
	DefuseCount          int `json:"defuse_count" gorm:"default:0"`
	SnipeCount           int `json:"snipe_count" gorm:"default:0"`
	FirstBloodCount      int `json:"first_blood_count" gorm:"default:0"`
	FirstDeathCount      int `json:"first_death_count" gorm:"default:0"`
	ZoneTierCaptureCount int `json:"zone_tier_capture_count" gorm:"default:0"`
	CtlAttackRounds      int `json:"ctl_attack_rounds" gorm:"default:0"`
	CtlDefenseRounds     int `json:"ctl_defense_rounds" gorm:"default:0"`

	NonTradedKills  int    `json:"non_traded_kills" gorm:"default:0"`
	HighestStreak   int    `json:"highest_streak" gorm:"default:0"`
	DataQualityNote string `json:"data_quality_note" gorm:"size:200"`
	Source          string `json:"source" gorm:"size:50"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Match  Match  `json:"match" gorm:"foreignKey:MatchID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
}

func (PlayerMapStats) TableName() string { return "player_map_stats" }

type PlayerMatchStats struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	MatchID      uint      `json:"match_id" gorm:"uniqueIndex:idx_player_match_stat_unique"`
	PlayerID     uint      `json:"player_id" gorm:"uniqueIndex:idx_player_match_stat_unique;index"`
	TeamID       uint      `json:"team_id" gorm:"index"`
	MapsPlayed   int       `json:"maps_played" gorm:"default:0"`
	TotalKills   int       `json:"total_kills" gorm:"default:0"`
	TotalDeaths  int       `json:"total_deaths" gorm:"default:0"`
	TotalAssists int       `json:"total_assists" gorm:"default:0"`
	TotalDamage  int       `json:"total_damage" gorm:"default:0"`
	KDRatio      float64   `json:"kd_ratio" gorm:"type:decimal(4,2);default:0"`
	KDARatio     float64   `json:"kda_ratio" gorm:"type:decimal(4,2);default:0"`
	ADR          float64   `json:"adr" gorm:"type:decimal(6,2);default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Match  Match  `json:"match" gorm:"foreignKey:MatchID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
}

func (PlayerMatchStats) TableName() string { return "player_match_stats" }

type PlayerTournamentStats struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	PlayerID     uint    `json:"player_id" gorm:"uniqueIndex:idx_player_tournament_stat_unique"`
	TeamID       uint    `json:"team_id" gorm:"index"`
	TournamentID uint    `json:"tournament_id" gorm:"uniqueIndex:idx_player_tournament_stat_unique;index"`
	TotalKills   int     `json:"total_kills"`
	TotalDeaths  int     `json:"total_deaths"`
	TotalAssists int     `json:"total_assists"`
	TotalDamage  int     `json:"total_damage"`
	KDRatio      float64 `json:"kd_ratio"`
	KDARatio     float64 `json:"kda_ratio"`

	Rank             *int    `json:"rank"`
	OverallPlusMinus int     `json:"overall_plus_minus" gorm:"default:0"`
	OverallMaps      int     `json:"overall_maps" gorm:"default:0"`
	SndKills         int     `json:"snd_kills" gorm:"default:0"`
	SndDeaths        int     `json:"snd_deaths" gorm:"default:0"`
	SndKDRatio       float64 `json:"snd_kd_ratio" gorm:"default:0"`
	SndPlusMinus     int     `json:"snd_plus_minus" gorm:"default:0"`
	SndKPerMap       float64 `json:"snd_k_per_map" gorm:"default:0"`
	SndFirstKills    int     `json:"snd_first_kills" gorm:"default:0"`
	SndMaps          int     `json:"snd_maps" gorm:"default:0"`
	HpKills          int     `json:"hp_kills" gorm:"default:0"`
	HpDeaths         int     `json:"hp_deaths" gorm:"default:0"`
	HpKDRatio        float64 `json:"hp_kd_ratio" gorm:"default:0"`
	HpPlusMinus      int     `json:"hp_plus_minus" gorm:"default:0"`
	HpKPerMap        float64 `json:"hp_k_per_map" gorm:"default:0"`
	HpTimeMilliseconds int   `json:"hp_time_milliseconds" gorm:"default:0"`
	HpMaps           int     `json:"hp_maps" gorm:"default:0"`
	ControlKills     int     `json:"control_kills" gorm:"default:0"`
	ControlDeaths    int     `json:"control_deaths" gorm:"default:0"`
	ControlKDRatio   float64 `json:"control_kd_ratio" gorm:"default:0"`
	ControlPlusMinus int     `json:"control_plus_minus" gorm:"default:0"`
	ControlKPerMap   float64 `json:"control_k_per_map" gorm:"default:0"`
	ControlCaptures  int     `json:"control_captures" gorm:"default:0"`
	ControlMaps      int     `json:"control_maps" gorm:"default:0"`

	Player     Player     `json:"player" gorm:"foreignKey:PlayerID"`
	Team       Team       `json:"team" gorm:"foreignKey:TeamID"`
	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
}

func (PlayerTournamentStats) TableName() string { return "player_tournament_stats" }

type TeamTournamentStats struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	TournamentID  uint      `json:"tournament_id" gorm:"index"`
	TeamID        uint      `json:"team_id" gorm:"index"`
	Placement     *int      `json:"placement"`
	MatchesPlayed int       `json:"matches_played" gorm:"default:0"`
	MatchesWon    int       `json:"matches_won" gorm:"default:0"`
	MatchesLost   int       `json:"matches_lost" gorm:"default:0"`
	MapsPlayed    int       `json:"maps_played" gorm:"default:0"`
	MapsWon       int       `json:"maps_won" gorm:"default:0"`
	MapsLost      int       `json:"maps_lost" gorm:"default:0"`
	PrizeMoney    float64   `json:"prize_money" gorm:"type:decimal(10,2);default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Tournament Tournament `json:"tournament" gorm:"foreignKey:TournamentID"`
	Team       Team       `json:"team" gorm:"foreignKey:TeamID"`
}

func (TeamTournamentStats) TableName() string { return "team_tournament_stats" }
