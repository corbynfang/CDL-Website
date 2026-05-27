package models

import "time"

type Team struct {
	ID             uint       `json:"id" gorm:"primaryKey;column:id"`
	Name           string     `json:"name" gorm:"not null;size:200;column:name"`
	Abbreviation   string     `json:"abbreviation" gorm:"not null;size:10;column:abbreviation"`
	City           string     `json:"city,omitempty" gorm:"size:100;column:city"`
	LogoURL        string     `json:"logo_url,omitempty" gorm:"column:logo_url"`
	PrimaryColor   string     `json:"primary_color" gorm:"size:7;column:primary_color"`
	SecondaryColor string     `json:"secondary_color" gorm:"size:7;column:secondary_color"`
	FoundedDate    *time.Time `json:"founded_date" gorm:"column:founded_date"`
	IsActive       bool       `json:"is_active" gorm:"default:true;column:is_active"`

	FranchiseID *uint  `json:"franchise_id" gorm:"column:franchise_id;index"`
	GameCode    string `json:"game_code" gorm:"size:10;column:game_code"`

	IsCDLFranchise     bool   `json:"is_cdl_franchise" gorm:"default:false;column:is_cdl_franchise;index"`
	TeamClassification string `json:"team_classification" gorm:"size:60;column:team_classification"`
	DoNotMerge         bool   `json:"do_not_merge" gorm:"default:false;column:do_not_merge"`

	ValidFrom *time.Time `json:"valid_from" gorm:"column:valid_from"`
	ValidTo   *time.Time `json:"valid_to" gorm:"column:valid_to"`

	NeedsManualReview bool   `json:"needs_manual_review" gorm:"default:false;column:needs_manual_review"`
	Source            string `json:"source" gorm:"size:100;column:source"`

	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Franchise *Franchise `json:"franchise" gorm:"foreignKey:FranchiseID"`
	Players   []Player   `json:"players" gorm:"many2many:team_rosters;"`
}

func (Team) TableName() string { return "teams" }

type TeamRoster struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	TeamID    uint       `json:"team_id" gorm:"index"`
	PlayerID  uint       `json:"player_id" gorm:"index"`
	SeasonID  uint       `json:"season_id" gorm:"index"`
	Role      string     `json:"role" gorm:"size:50"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsStarter bool       `json:"is_starter" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	Team   Team   `json:"team" gorm:"foreignKey:TeamID"`
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
	Season Season `json:"season" gorm:"foreignKey:SeasonID"`
}

type Coach struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null;size:100"`
	TeamID   uint   `json:"team_id"`
	SeasonID uint   `json:"season_id"`
}

func (Coach) TableName() string { return "coaches" }
