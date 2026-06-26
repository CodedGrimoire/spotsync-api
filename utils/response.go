package utils

import "github.com/labstack/echo/v4"

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, map[string]interface{}{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}
