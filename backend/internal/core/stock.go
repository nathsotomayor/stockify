package core

import (
	"time"

	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	Ticker     string    `gorm:"not null;index" json:"ticker"`
	Company    string    `gorm:"not null" json:"company"`
	Brokerage  string    `gorm:"not null" json:"brokerage"`
	Action     string    `gorm:"not null" json:"action"`
	RatingTo   string    `gorm:"not null" json:"rating_to"`
	RatingFrom *string   `json:"rating_from,omitempty"`
	TargetTo   *float64  `gorm:"type:decimal(10,2)" json:"target_to,omitempty"`
	TargetFrom *float64  `gorm:"type:decimal(10,2)" json:"target_from,omitempty"`
	Time       time.Time `json:"time"`
}
