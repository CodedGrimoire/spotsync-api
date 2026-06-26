package models

import "time"

type Reservation struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `json:"user_id"`
	ZoneID       uint        `json:"zone_id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `gorm:"default:active" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	User         User        `gorm:"foreignKey:UserID" json:"user"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID" json:"zone"`
}
