package models

import "time"

type Tournament struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	SeasonID         uint       `json:"season_id" gorm:"index"`
	Name             string     `json:"name" gorm:"not null;size:200"`
	Slug             string     `json:"slug" gorm:"size:200;index"`
	TournamentType   string     `json:"tournament_type" gorm:"size:50"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	PrizePool        *float64   `json:"prize_pool" gorm:"type:decimal(12,2)"`
	Location         string     `json:"location" gorm:"size:200"`
	Country          string     `json:"country" gorm:"size:3"`
	IsLAN            bool       `json:"is_lan" gorm:"default:false"`
	LogoURL          string     `json:"logo_url" gorm:"size:500"`
	TournamentFormat string     `json:"tournament_format" gorm:"size:50"`
	SourceEventURL   string     `json:"source_event_url,omitempty" gorm:"column:liquipedia_url"`
	SourceURL        string     `json:"-" gorm:"column:breaking_point_url"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Season Season `json:"season" gorm:"foreignKey:SeasonID"`
}

func (Tournament) TableName() string { return "tournaments" }
