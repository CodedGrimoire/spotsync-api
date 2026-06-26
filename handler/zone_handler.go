package handler

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateParkingZoneRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zoneResponse, err := h.zoneService.CreateZone(req)
	if errors.Is(err, service.ErrInvalidZoneType) {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone type", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create parking zone", nil)
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Parking zone created successfully", zoneResponse)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zoneResponses, err := h.zoneService.GetAllZones()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve parking zones", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zones retrieved successfully", zoneResponses)
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	id, err := parseZoneID(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	zoneResponse, err := h.zoneService.GetZoneByID(id)
	if errors.Is(err, service.ErrZoneNotFound) {
		return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve parking zone", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zone retrieved successfully", zoneResponse)
}

func (h *ZoneHandler) UpdateZone(c echo.Context) error {
	id, err := parseZoneID(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	var req dto.UpdateParkingZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zoneResponse, err := h.zoneService.UpdateZone(id, req)
	if errors.Is(err, service.ErrZoneNotFound) {
		return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
	}
	if errors.Is(err, service.ErrInvalidZoneType) {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone type", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update parking zone", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zone updated successfully", zoneResponse)
}

func (h *ZoneHandler) DeleteZone(c echo.Context) error {
	id, err := parseZoneID(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid zone ID", nil)
	}

	if err := h.zoneService.DeleteZone(id); errors.Is(err, service.ErrZoneNotFound) {
		return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
	} else if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete parking zone", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Parking zone deleted successfully", nil)
}

func parseZoneID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 0)
	if err != nil || id == 0 {
		if err == nil {
			err = errors.New("invalid zone ID")
		}

		return 0, err
	}

	return uint(id), nil
}
