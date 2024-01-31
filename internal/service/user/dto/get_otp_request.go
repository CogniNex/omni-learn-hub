package dto

type UserGetOtpRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,max=12,min=9"`
}
