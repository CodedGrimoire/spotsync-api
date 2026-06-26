package routes

import (
	"log"
	"os"

	"spotsync-api/handler"
	appMiddleware "spotsync-api/middleware"
	"spotsync-api/repository"
	"spotsync-api/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api/v1")
	auth := api.Group("/auth")
	protected := api.Group("")
	protected.Use(appMiddleware.JWTMiddleware())
	admin := api.Group("")
	admin.Use(appMiddleware.JWTMiddleware())
	admin.Use(appMiddleware.RequireRole("admin"))

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)
	zoneRepo := repository.NewZoneRepository(db)
	zoneService := service.NewZoneService(zoneRepo)
	zoneHandler := handler.NewZoneHandler(zoneService)
	reservationRepo := repository.NewReservationRepository(db)
	reservationService := service.NewReservationService(reservationRepo)
	reservationHandler := handler.NewReservationHandler(reservationService)

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
