package otp

import (
	crytporand "crypto/rand"
	"encoding/base32"
)

type Generator interface {
	RandomSecret(length int) string
}

type GOTPGenerator struct{}

func NewGOTPGenerator() *GOTPGenerator {
	return &GOTPGenerator{}
}

func (g *GOTPGenerator) RandomSecret(length int) string {
	var result string
	secret := make([]byte, length)
	gen, err := crytporand.Read(secret)
	if err != nil || gen != length {
		// error reading random, return empty string
		return result
	}
	var encoder = base32.StdEncoding.WithPadding(base32.NoPadding)
	result = encoder.EncodeToString(secret)
	return result
}
