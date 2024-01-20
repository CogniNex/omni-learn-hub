package service

import service "omni-learn-hub/internal/service/userservice"

type Services struct {
	Users Users
}

type Users interface {
}

type Deps struct {
}

func NewServices(deps Deps) *Services {
	userService := service.NewUserService()
	return &Services{
		Users: userService,
	}

}
