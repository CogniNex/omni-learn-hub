package dto

type UserSignUpInput struct {
	PhoneNumber          string `json:"phone_number" binding:"required,max=13"`
	Password             string `json:"password" binding:"required,min=8,max=64"`
	PasswordVerification string `json:"password_verification" binding:"required,min=8,max=64"`
}
