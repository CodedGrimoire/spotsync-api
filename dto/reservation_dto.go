package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationZoneResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ReservationUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type ReservationResponse struct {
	ID           uint                     `json:"id"`
	UserID       uint                     `json:"user_id"`
	ZoneID       uint                     `json:"zone_id"`
	LicensePlate string                   `json:"license_plate"`
	Status       string                   `json:"status"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	Zone         *ReservationZoneResponse `json:"zone,omitempty"`
	User         *ReservationUserResponse `json:"user,omitempty"`
}
