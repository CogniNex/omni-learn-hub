package request

type UserGetOtpRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,max=12,min=9"`
}
