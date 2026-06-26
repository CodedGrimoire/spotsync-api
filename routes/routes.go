package routes

import (
	"log"
	"os"

	"spotsync-api/handler"
	"spotsync-api/repository"
	"spotsync-api/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	api := e.Group("/api/v1")
	auth := api.Group("/auth")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(authService)

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
