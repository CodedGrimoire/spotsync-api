package routes

import (
	"spotsync-api/handler"
	appMiddleware "spotsync-api/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	zoneHandler *handler.ZoneHandler,
	reservationHandler *handler.ReservationHandler,
) {
	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	protected := api.Group("")
	protected.Use(appMiddleware.JWTMiddleware())
	admin := api.Group("")
	admin.Use(appMiddleware.JWTMiddleware())
	admin.Use(appMiddleware.RequireRole("admin"))

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	api.GET("/zones", zoneHandler.GetAllZones)
	api.GET("/zones/:id", zoneHandler.GetZoneByID)

	admin.POST("/zones", zoneHandler.CreateZone)
	admin.PUT("/zones/:id", zoneHandler.UpdateZone)
	admin.DELETE("/zones/:id", zoneHandler.DeleteZone)

	protected.POST("/reservations", reservationHandler.CreateReservation)
	protected.GET("/reservations/my-reservations", reservationHandler.GetMyReservations)
	protected.DELETE("/reservations/:id", reservationHandler.CancelReservation)

	admin.GET("/reservations", reservationHandler.GetAllReservations)
}
