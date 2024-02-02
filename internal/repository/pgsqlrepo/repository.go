package pgsqlrepo

import (
	"omni-learn-hub/internal/repository"
	"omni-learn-hub/internal/repository/pgsqlrepo/otp_code"
	"omni-learn-hub/internal/repository/pgsqlrepo/user"
	"omni-learn-hub/pkg/postgres"
)

type Repositories struct {
	Users    repository.Users
	OtpCodes repository.OtpCodes
}

func NewRepositories(db *postgres.Postgres) *Repositories {
	return &Repositories{
		Users:    user.NewUsersRepo(db),
		OtpCodes: otp_code.NewOtpCodesRepo(db),
	}
}
