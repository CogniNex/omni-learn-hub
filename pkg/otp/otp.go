package otp

import (
	"math/rand"
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
	characters := "0123456789"
	otp := make([]byte, g.Length)

	for i := range otp {
		otp[i] = characters[rand.Intn(len(characters))]
	}

	return string(otp)
}
