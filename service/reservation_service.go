package service

import (
	"errors"

	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"

	"gorm.io/gorm"
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.ReservationResponse, error)
	GetAllReservations() ([]dto.ReservationResponse, error)
	CancelReservation(userID uint, role string, reservationID uint) error
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationService{reservationRepo: reservationRepo}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation, err := s.reservationRepo.CreateWithTransaction(userID, req.ZoneID, req.LicensePlate)
	if errors.Is(err, repository.ErrZoneFull) {
		return nil, ErrZoneFull
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrZoneNotFound
	}
	if err != nil {
		return nil, err
	}

	response := toReservationResponse(reservation)

	return &response, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReservationResponse, 0, len(reservations))
	for i := range reservations {
		responses = append(responses, toReservationResponse(&reservations[i]))
	}

	return responses, nil
}

func (s *reservationService) GetAllReservations() ([]dto.ReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ReservationResponse, 0, len(reservations))
	for i := range reservations {
		responses = append(responses, toReservationResponse(&reservations[i]))
	}

	return responses, nil
}

func (s *reservationService) CancelReservation(userID uint, role string, reservationID uint) error {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrReservationNotFound
	}
	if err != nil {
		return err
	}

	if role != "admin" && reservation.UserID != userID {
		return ErrReservationForbidden
	}

	if reservation.Status == "cancelled" {
		return ErrReservationAlreadyCancelled
	}

	if reservation.Status == "completed" {
		return ErrReservationAlreadyCompleted
	}

	reservation.Status = "cancelled"

	return s.reservationRepo.Update(reservation)
}

func toReservationResponse(reservation *models.Reservation) dto.ReservationResponse {
	response := dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}

	if reservation.Zone.ID != 0 {
		response.Zone = &dto.ReservationZoneResponse{
			ID:   reservation.Zone.ID,
			Name: reservation.Zone.Name,
			Type: reservation.Zone.Type,
		}
	}

	if reservation.User.ID != 0 {
		response.User = &dto.ReservationUserResponse{
			ID:    reservation.User.ID,
			Name:  reservation.User.Name,
			Email: reservation.User.Email,
			Role:  reservation.User.Role,
		}
	}

	return response
}
