package main

import (
	"log"
	"os"

	"spotsync-api/config"
	"spotsync-api/models"
	"spotsync-api/routes"
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
	routes.RegisterRoutes(e, db)

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
