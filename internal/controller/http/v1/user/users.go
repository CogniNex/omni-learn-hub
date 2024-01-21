package user

import (
	"github.com/gin-gonic/gin"
	"net/http"

	userService "omni-learn-hub/internal/service/user"
	"omni-learn-hub/pkg/logger"
)

type userRoutes struct {
	userService userService.Users
	logger      logger.Interface
}

func NewUserRoutes(handler *gin.RouterGroup, userService userService.Users, logger logger.Interface) {
	r := &userRoutes{userService: userService, logger: logger}
	h := handler.Group("/user")
	{
		h.GET("/joke", r.joke)
	}
}

func (r *userRoutes) joke(c *gin.Context) {
	c.JSON(http.StatusOK, "Аслану подходит прическа ")
}
