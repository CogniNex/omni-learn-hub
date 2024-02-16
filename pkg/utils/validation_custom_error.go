package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

type ApiValidationError struct {
	Field string `json:"field"`
	Msg   string `json:"errorMessage"`
}

func ValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {

		out := make([]ApiValidationError, len(ve))
		for i, fe := range ve {
			out[i] = ApiValidationError{fe.Field(), msgForTag(fe.Tag())}
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": out, "success": false})
	}
	return
}
