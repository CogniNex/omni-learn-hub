package service

import (
	"omni-learn-hub/internal/repository/pgsqlrepo"
	userService "omni-learn-hub/internal/service/user"
)

type Services struct {
	Users userService.Users
}

type Deps struct {
	Repos *pgsqlrepo.Repositories
}

func NewServices(deps Deps) *Services {
	userService := userService.NewUserService()
	return &Services{
		Users: userService,
	}

}
