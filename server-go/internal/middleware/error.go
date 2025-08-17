package middleware

import (
	"net/http"

	"server/internal/repository"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		status := http.StatusInternalServerError
		appErr := AppError{Code: "INTERNAL", Message: err.Error()}

		switch err {
		case repository.ErrNotFound:
			status = http.StatusNotFound
			appErr = AppError{Code: "NOT_FOUND", Message: "resource not found"}
		default:
			// có thể phân loại thêm theo type/assert
		}

		c.AbortWithStatusJSON(status, appErr)
	}
}
