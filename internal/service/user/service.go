package user

import "omni-learn-hub/internal/repository/pgsqlrepo/user"

type UsersService struct {
	repo user.UsersRepo
}

type Users interface {
}

func NewUserService() *UsersService {
	return &UsersService{}
}
