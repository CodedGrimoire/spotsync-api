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

type ReservationHandler struct {
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	reservationResponse, err := h.reservationService.CreateReservation(userID, req)
	if errors.Is(err, service.ErrZoneNotFound) {
		return utils.ErrorResponse(c, http.StatusNotFound, "Parking zone not found", nil)
	}
	if errors.Is(err, service.ErrZoneFull) {
		return utils.ErrorResponse(c, http.StatusConflict, "Parking zone is full", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create reservation", nil)
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Reservation confirmed successfully", reservationResponse)
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	reservationResponses, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reservations", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "My reservations retrieved successfully", reservationResponses)
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservationResponses, err := h.reservationService.GetAllReservations()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reservations", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Reservations retrieved successfully", reservationResponses)
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	role := getRoleFromContext(c)
	reservationID, err := parseReservationID(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reservation ID", nil)
	}

	err = h.reservationService.CancelReservation(userID, role, reservationID)
	if errors.Is(err, service.ErrReservationNotFound) {
		return utils.ErrorResponse(c, http.StatusNotFound, "Reservation not found", nil)
	}
	if errors.Is(err, service.ErrReservationForbidden) {
		return utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", nil)
	}
	if errors.Is(err, service.ErrReservationAlreadyCancelled) {
		return utils.ErrorResponse(c, http.StatusConflict, "Reservation is already cancelled", nil)
	}
	if errors.Is(err, service.ErrReservationAlreadyCompleted) {
		return utils.ErrorResponse(c, http.StatusConflict, "Completed reservation cannot be cancelled", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel reservation", nil)
	}

	return utils.SuccessResponse(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func getUserIDFromContext(c echo.Context) (uint, bool) {
	userID, ok := c.Get("user_id").(uint)
	if !ok || userID == 0 {
		return 0, false
	}

	return userID, true
}

func getRoleFromContext(c echo.Context) string {
	role, _ := c.Get("role").(string)

	return role
}

func parseReservationID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 0)
	if err != nil || id == 0 {
		if err == nil {
			err = errors.New("invalid reservation ID")
		}

		return 0, err
	}

	return uint(id), nil
}
