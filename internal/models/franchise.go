package models

import "time"

// Franchise is a stable CDL franchise slot that persists across team rebrands.
// Minnesota RØKKR and G2 Minnesota are separate Team rows both pointing to the same Franchise.
type Franchise struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FranchiseKey string    `json:"franchise_key" gorm:"uniqueIndex;not null;size:100"`
	Name         string    `json:"name" gorm:"not null;size:200"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Franchise) TableName() string { return "franchises" }
