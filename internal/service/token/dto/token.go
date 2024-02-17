package dto

import "github.com/gofrs/uuid"

type TokenDto struct {
	UserId      uuid.UUID
	FirstName   string
	LastName    string
	PhoneNumber string
	Roles       []string
}
