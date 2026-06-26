package models

import "time"

type ParkingZone struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Type          string    `gorm:"size:30;not null;check:type IN ('general','ev_charging','covered')" json:"type"`
	TotalCapacity int       `gorm:"not null;check:total_capacity > 0" json:"total_capacity"`
	PricePerHour  float64   `gorm:"not null;check:price_per_hour > 0" json:"price_per_hour"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
