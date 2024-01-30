package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"omni-learn-hub/internal/service/user/dto"

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
		h.POST("/sign-up", r.signUp)
	}
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body dto.UserSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} string "ok"
// @Failure 500 {object} string "ok"
// @Failure default {object} string "ok"
// @Router /users/sign-up [post]
func (r *userRoutes) signUp(c *gin.Context) {
	var inp dto.UserSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")

		return
	}

	if err := r.userService.SignUp(c.Request.Context(), dto.UserSignUpInput{
		PhoneNumber:          inp.PhoneNumber,
		Password:             inp.Password,
		PasswordVerification: inp.PasswordVerification,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())

		return
	}

	c.Status(http.StatusCreated)
}
