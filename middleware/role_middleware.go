package middleware

import (
	"net/http"

	"spotsync-api/utils"

	"github.com/labstack/echo/v4"
)

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok || role == "" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
			}

			if isRoleAllowed(role, allowedRoles) {
				return next(c)
			}

			return utils.ErrorResponse(c, http.StatusForbidden, "Forbidden: insufficient permissions", nil)
		}
	}
}

func isRoleAllowed(role string, allowedRoles []string) bool {
	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}

	return false
}
