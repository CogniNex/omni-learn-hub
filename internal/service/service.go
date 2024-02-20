package service

import (
	"omni-learn-hub/config"
	"omni-learn-hub/internal/repository/pgsqlrepo"
	"omni-learn-hub/internal/service/token"
	userService "omni-learn-hub/internal/service/user"
	"omni-learn-hub/pkg/hash"
	"omni-learn-hub/pkg/otp"
	"omni-learn-hub/pkg/sms"
	"omni-learn-hub/pkg/whatsapp"
)

type Services struct {
	Users userService.Users
	Token token.Tokens
}

type Deps struct {
	Repos    *pgsqlrepo.Repositories
	Hasher   hash.PasswordHasher
	Otp      otp.Generator
	SMS      sms.SMSClient
	Whatsapp whatsapp.WhatsappClient
	Cfg      *config.Config
}

func NewServices(deps Deps) *Services {
	t := token.NewTokenService(deps.Cfg)
	u := userService.NewUserService(deps.Repos.Users, deps.Repos.OtpCodes, deps.Hasher, deps.Otp, deps.SMS, *t, deps.Whatsapp)
	return &Services{
		Users: u,
		Token: t,
	}

}
