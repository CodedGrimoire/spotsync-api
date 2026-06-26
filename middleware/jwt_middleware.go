package middleware

import (
	"net/http"
	"os"
	"strings"

	"spotsync-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				return utils.ErrorResponse(c, http.StatusInternalServerError, "JWT secret is not configured", nil)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Missing authorization header", nil)
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header", nil)
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			}

			userID, ok := getUserIDFromClaims(claims)
			if !ok {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			}

			role, ok := claims["role"].(string)
			if !ok || role == "" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			}

			c.Set("user_id", userID)
			c.Set("role", role)

			return next(c)
		}
	}
}

func getUserIDFromClaims(claims jwt.MapClaims) (uint, bool) {
	userIDValue, ok := claims["user_id"].(float64)
	maxUint := ^uint(0)
	if !ok || userIDValue <= 0 || userIDValue > float64(maxUint) || userIDValue != float64(uint(userIDValue)) {
		return 0, false
	}

	return uint(userIDValue), true
}
