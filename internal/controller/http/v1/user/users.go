package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"omni-learn-hub/internal/service/user/dto/request"

	userService "omni-learn-hub/internal/service/user"
	"omni-learn-hub/pkg/logger"
)

type userRoutes struct {
	userService userService.Users
	logger      logger.Interface
}

func NewUserRoutes(handler *gin.RouterGroup, userService userService.Users, logger logger.Interface) {
	r := &userRoutes{userService: userService, logger: logger}
	h := handler.Group("/users")
	{
		h.POST("/sign-up", r.signUp)
		h.POST("/get-otp", r.getOtp)
	}
}

// @Summary User SignUp
// @Tags Users
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body dto.UserSignUpRequest true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} string "ok"
// @Failure 500 {object} string "ok"
// @Failure default {object} string "ok"
// @Router /api/v1/users/sign-up [post]
func (r *userRoutes) signUp(c *gin.Context) {
	var req request.UserSignUpRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")

		return
	}

	response := r.userService.SignUp(c.Request.Context(), request.UserSignUpRequest{
		PhoneNumber:          req.PhoneNumber,
		OtpCode:              req.OtpCode,
		Password:             req.Password,
		PasswordVerification: req.PasswordVerification,
	})

	if response.Success == false {
		c.AbortWithStatusJSON(http.StatusBadRequest, response)

		return
	}

	c.Status(http.StatusCreated)
	c.JSON(http.StatusCreated, response)

}

// @Summary Get Otp Code
// @Tags Users
// @Description user gets otp which need for registration
// @ModuleID userGetOtp
// @Accept  json
// @Produce  json
// @Param input body dto.UserGetOtpRequest true "get otp info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} string "ok"
// @Failure 500 {object} string "ok"
// @Failure default {object} string "ok"
// @Router /api/v1/users/get-otp [post]
func (r *userRoutes) getOtp(c *gin.Context) {

	var inp request.UserGetOtpRequest
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid input body")

		return
	}

	if err := r.userService.GetOtp(c.Request.Context(), request.UserGetOtpRequest{
		PhoneNumber: inp.PhoneNumber,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())

		return
	}

	c.Status(http.StatusOK)

}
