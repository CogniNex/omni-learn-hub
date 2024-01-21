package pgsqlrepo

import (
	"omni-learn-hub/internal/repository"
	"omni-learn-hub/internal/repository/pgsqlrepo/user"
	"omni-learn-hub/pkg/postgres"
)

type Repositories struct {
	Users repository.Users
}

func NewRepositories(db *postgres.Postgres) *Repositories {
	return &Repositories{
		Users: user.NewUsersRepo(db),
	}
}
