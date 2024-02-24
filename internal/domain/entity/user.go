package entity

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	UserID               uuid.UUID `db:"user_id"`
	PhoneNumber          string    `db:"phone_number"`
	PasswordHash         string    `db:"password_hash"`
	PasswordSalt         string    `db:"password_salt"`
	RefreshToken         string    `db:"refresh_token"`
	RefreshTokenExpireIn time.Time `db:"refresh_expires_in"`
}
