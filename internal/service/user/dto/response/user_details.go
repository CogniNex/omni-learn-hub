package response

import (
	"github.com/gofrs/uuid"
	"time"
)

type UserDetails struct {
	UserID               uuid.UUID `db:"user_id"`
	PhoneNumber          string    `db:"phone_number"`
	PasswordHash         string    `db:"password_hash"`
	PasswordSalt         string    `db:"password_salt"`
	FirstName            string    `db:"first_name"`
	EntityID             int       `db:"entity_id"`
	EntityTypeID         int       `db:"entity_type_id"`
	Lastname             string    `db:"last_name"`
	DateOfBirth          time.Time `db:"date_of_birth"`
	LanguageID           int       `db:"language_id"`
	Email                string    `db:"email"`
	IsActive             bool      `db:"is_active"`
	RefreshToken         string    `db:"refresh_token"`
	RefreshTokenExpireIn time.Time `db:"refresh_expires_in"`
}
