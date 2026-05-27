package models

import "time"

type Season struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"not null;size:100"`
	GameTitle string     `json:"game_title" gorm:"not null;size:100"`
	GameCode  string     `json:"game_code" gorm:"size:10"` // BO6, CW, MW2, MW3, VG
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsActive  bool       `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (Season) TableName() string { return "seasons" }
