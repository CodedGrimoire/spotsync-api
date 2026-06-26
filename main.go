package main

import (
	"log"
	"os"

	"spotsync-api/config"
	"spotsync-api/handler"
	"spotsync-api/models"
	"spotsync-api/repository"
	"spotsync-api/routes"
	"spotsync-api/service"
	"spotsync-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db := config.ConnectDB()

	if err := db.AutoMigrate(
		&models.User{},
		&models.ParkingZone{},
		&models.Reservation{},
	); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

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

	routes.RegisterRoutes(e, authHandler, zoneHandler, reservationHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"success": true,
			"message": "SpotSync API is running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
