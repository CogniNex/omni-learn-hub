package otp

import (
	crytporand "crypto/rand"
	"encoding/base32"
)

type Generator interface {
	RandomSecret() string
}

type GOTPGenerator struct {
	Length int
}

func NewGOTPGenerator(length int) *GOTPGenerator {
	return &GOTPGenerator{Length: length}
}

func (g *GOTPGenerator) RandomSecret() string {
	var result string
	secret := make([]byte, g.Length)
	gen, err := crytporand.Read(secret)
	if err != nil || gen != g.Length {
		// error reading random, return empty string
		return result
	}
	var encoder = base32.StdEncoding.WithPadding(base32.NoPadding)
	result = encoder.EncodeToString(secret)
	return result
}
