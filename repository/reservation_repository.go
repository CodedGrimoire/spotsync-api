package repository

import (
	"errors"

	"spotsync-api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrZoneFull = errors.New("parking zone is full")

type ReservationRepository interface {
	CreateWithTransaction(userID uint, reqZoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindAll() ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	Update(reservation *models.Reservation) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateWithTransaction(userID uint, reqZoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation *models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, reqZoneID).Error; err != nil {
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", reqZoneID, "active").
			Count(&activeCount).
			Error; err != nil {
			return err
		}

		if activeCount >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		reservation = &models.Reservation{
			UserID:       userID,
			ZoneID:       reqZoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(reservation).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reservations).
		Error

	return reservations, err
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("User").
		Preload("Zone").
		Order("created_at DESC").
		Find(&reservations).
		Error

	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := r.db.Preload("User").Preload("Zone").First(&reservation, id).Error; err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}
