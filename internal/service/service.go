package service

import (
	"omni-learn-hub/internal/repository/pgsqlrepo"
	userService "omni-learn-hub/internal/service/user"
	"omni-learn-hub/pkg/hash"
	"omni-learn-hub/pkg/otp"
	"omni-learn-hub/pkg/sms"
)

type Services struct {
	Users userService.Users
}

type Deps struct {
	Repos  *pgsqlrepo.Repositories
	Hasher hash.PasswordHasher
	Otp    otp.Generator
	SMS    sms.SMSClient
}

func NewServices(deps Deps) *Services {
	u := userService.NewUserService(deps.Repos.Users, deps.Hasher, deps.Otp, deps.SMS)
	return &Services{
		Users: u,
	}

}
