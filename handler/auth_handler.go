package handler

import (
	"errors"

	"spotsync-api/dto"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Validation failed", err.Error())
	}

	userResponse, err := h.authService.Register(req)
	if errors.Is(err, service.ErrEmailAlreadyExists) {
		return utils.ErrorResponse(c, 400, "Email already exists", nil)
	}
	if errors.Is(err, service.ErrInvalidRole) {
		return utils.ErrorResponse(c, 400, "Invalid role", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to register user", nil)
	}

	return utils.SuccessResponse(c, 201, "User registered successfully", userResponse)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Validation failed", err.Error())
	}

	loginResponse, err := h.authService.Login(req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		return utils.ErrorResponse(c, 401, "Invalid email or password", nil)
	}
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to login", nil)
	}

	return utils.SuccessResponse(c, 200, "Login successful", loginResponse)
}
