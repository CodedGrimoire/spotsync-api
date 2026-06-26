package service

import (
	"errors"

	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"

	"gorm.io/gorm"
)

type ZoneService interface {
	CreateZone(req dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	GetAllZones() ([]dto.ParkingZoneResponse, error)
	GetZoneByID(id uint) (*dto.ParkingZoneResponse, error)
	UpdateZone(id uint, req dto.UpdateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	if !isValidZoneType(req.Type) {
		return nil, ErrInvalidZoneType
	}

	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}

	response := toParkingZoneResponse(zone, zone.TotalCapacity)

	return &response, nil
}

func (s *zoneService) GetAllZones() ([]dto.ParkingZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ParkingZoneResponse, 0, len(zones))
	for i := range zones {
		availableSpots, err := s.getAvailableSpots(&zones[i])
		if err != nil {
			return nil, err
		}

		responses = append(responses, toParkingZoneResponse(&zones[i], availableSpots))
	}

	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ParkingZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrZoneNotFound
	}
	if err != nil {
		return nil, err
	}

	availableSpots, err := s.getAvailableSpots(zone)
	if err != nil {
		return nil, err
	}

	response := toParkingZoneResponse(zone, availableSpots)

	return &response, nil
}

func (s *zoneService) UpdateZone(id uint, req dto.UpdateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrZoneNotFound
	}
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		if !isValidZoneType(req.Type) {
			return nil, ErrInvalidZoneType
		}
		zone.Type = req.Type
	}
	if req.TotalCapacity != 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour != 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}

	availableSpots, err := s.getAvailableSpots(zone)
	if err != nil {
		return nil, err
	}

	response := toParkingZoneResponse(zone, availableSpots)

	return &response, nil
}

func (s *zoneService) DeleteZone(id uint) error {
	if _, err := s.zoneRepo.FindByID(id); errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrZoneNotFound
	} else if err != nil {
		return err
	}

	return s.zoneRepo.Delete(id)
}

func (s *zoneService) getAvailableSpots(zone *models.ParkingZone) (int, error) {
	activeReservations, err := s.zoneRepo.CountActiveReservations(zone.ID)
	if err != nil {
		return 0, err
	}

	availableSpots := zone.TotalCapacity - int(activeReservations)
	if availableSpots < 0 {
		return 0, nil
	}

	return availableSpots, nil
}

func isValidZoneType(zoneType string) bool {
	return zoneType == "general" || zoneType == "ev_charging" || zoneType == "covered"
}

func toParkingZoneResponse(zone *models.ParkingZone, availableSpots int) dto.ParkingZoneResponse {
	if availableSpots < 0 {
		availableSpots = 0
	}

	return dto.ParkingZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: availableSpots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}
}
