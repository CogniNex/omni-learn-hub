package hash

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) (string, string, error)
	CheckPassword(password, hashedPassword, salt string) bool
}

type BcryptPasswordHasher struct{}

// NewBcryptPasswordHasher creates a new instance of bcryptPasswordHasher
func NewBcryptPasswordHasher() PasswordHasher {
	return &BcryptPasswordHasher{}
}

func (h *BcryptPasswordHasher) HashPassword(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return string(hashedPassword), salt, nil
}

func (h *BcryptPasswordHasher) CheckPassword(password, hashedPassword, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	return err == nil
}

func generateSalt() (string, error) {
	saltBytes := make([]byte, 16) // 16 bytes for a 128-bit salt
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", errors.New("failed to generate salt")
	}

	return hex.EncodeToString(saltBytes), nil
}
