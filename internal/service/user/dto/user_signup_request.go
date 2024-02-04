package dto

type UserSignUpRequest struct {
	PhoneNumber          string `json:"phoneNumber" binding:"required,max=13"`
	OtpCode              string `json:"otpCode" binding:"required"`
	Password             string `json:"password" binding:"required,min=8,max=64"`
	PasswordVerification string `json:"passwordVerification" binding:"required,min=8,max=64"`
}
