package service

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRole        = errors.New("invalid role")
	ErrUserNotFound       = errors.New("user not found")
	ErrZoneNotFound       = errors.New("parking zone not found")
	ErrInvalidZoneType    = errors.New("invalid zone type")

	ErrReservationNotFound         = errors.New("reservation not found")
	ErrReservationForbidden        = errors.New("you are not allowed to modify this reservation")
	ErrZoneFull                    = errors.New("parking zone is full")
	ErrReservationAlreadyCancelled = errors.New("reservation is already cancelled")
	ErrReservationAlreadyCompleted = errors.New("completed reservation cannot be cancelled")
)
