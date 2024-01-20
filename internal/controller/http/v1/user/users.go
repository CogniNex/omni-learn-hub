package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"omni-learn-hub/internal/service"
	"omni-learn-hub/pkg/logger"
)

type userRoutes struct {
	userService service.Users
	logger      logger.Interface
}

func NewUserRoutes(handler *gin.RouterGroup, userService service.Users, logger logger.Interface) {
	r := &userRoutes{userService: userService, logger: logger}
	h := handler.Group("/user")
	{
		h.GET("/joke", r.joke)
	}
}

func (r *userRoutes) joke(c *gin.Context) {
	c.JSON(http.StatusOK, "Аслану подходит прическа ")
}
