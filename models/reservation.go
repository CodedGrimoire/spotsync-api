package models

import "time"

type Reservation struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	ZoneID       uint        `gorm:"not null;index" json:"zone_id"`
	LicensePlate string      `gorm:"size:15;not null" json:"license_plate"`
	Status       string      `gorm:"size:20;not null;default:active;check:status IN ('active','completed','cancelled')" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"user"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"zone"`
}
