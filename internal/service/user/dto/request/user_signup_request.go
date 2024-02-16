package request

type UserSignUpRequest struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required,max=13"`
	OtpCode     string `json:"otpCode" binding:"required"`
	Password    string `json:"password" binding:"required,min=8,max=64"`
	RoleId      int    `json:"roleId" binding:"required"`
}
