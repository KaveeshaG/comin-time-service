// internal/middleware/error_handler.go
package middleware

import (
	"fmt"
	"net/http"

	"github.com/Axontik/comin-time-service/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Handle AppError
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.HTTPStatus, appErr)
				return
			}

			// Handle GORM errors
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, errors.NewNotFoundError("Resource not found"))
				return
			}

			// Handle validation errors
			if verr, ok := err.(validator.ValidationErrors); ok {
				details := translateValidationErrors(verr)
				c.JSON(http.StatusBadRequest, errors.NewValidationError(details))
				return
			}

			// Handle unknown errors
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError("An unexpected error occurred"))
		}
	}
}

func translateValidationErrors(errs validator.ValidationErrors) []map[string]string {
	var details []map[string]string

	for _, err := range errs {
		detail := map[string]string{
			"field":   err.Field(),
			"tag":     err.Tag(),
			"value":   err.Param(),
			"message": getValidationErrorMessage(err),
		}
		details = append(details, detail)
	}

	return details
}

func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must not be longer than %s characters", err.Param())
	default:
		return "Invalid value"
	}
}
