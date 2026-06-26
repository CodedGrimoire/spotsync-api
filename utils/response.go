package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	statusCode := http.StatusInternalServerError
	message := http.StatusText(http.StatusInternalServerError)

	if echoErr, ok := err.(*echo.HTTPError); ok {
		statusCode = echoErr.Code
		message = http.StatusText(statusCode)

		if statusCode < http.StatusInternalServerError {
			if msg, ok := echoErr.Message.(string); ok && msg != "" {
				message = msg
			}
		}
	}

	_ = ErrorResponse(c, statusCode, message, nil)
}
