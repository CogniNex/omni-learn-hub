package entity

import "time"

type UserProfile struct {
	UserID       string    `db:"user_id"`
	Name         string    `db:"name"`
	EntityID     int       `db:"entity_id"`
	EntityTypeID int       `db:"entity_type_id"`
	Surname      string    `db:"surname"`
	DateOfBirth  time.Time `db:"date_of_birth"`
	LanguageID   int       `db:"language_id"`
	Email        string    `db:"email"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
