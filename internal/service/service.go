package service

import (
	"omni-learn-hub/internal/repository/pgsqlrepo"
	userService "omni-learn-hub/internal/service/user"
	"omni-learn-hub/pkg/hash"
)

type Services struct {
	Users userService.Users
}

type Deps struct {
	Repos  *pgsqlrepo.Repositories
	Hasher hash.PasswordHasher
}

func NewServices(deps Deps) *Services {
	u := userService.NewUserService(deps.Repos.Users, deps.Hasher)
	return &Services{
		Users: u,
	}

}
