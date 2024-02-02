package entity

import "time"

type OtpCode struct {
	OtpID              int       `db:"otp_id"`
	PhoneNumber        string    `db:"phone_number"`
	Code               string    `db:"code"`
	IsVerified         bool      `db:"is_verified"`
	GenerationAttempts int       `db:"generation_attempts"`
	CreatedAt          time.Time `db:"created_at"`
	ExpiresAt          time.Time `db:"expires_at"`
}
