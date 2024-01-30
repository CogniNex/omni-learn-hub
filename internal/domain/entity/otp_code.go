package entity

import "time"

type OtpCode struct {
	OtpID      string    `db:"otp_id"`
	UserID     string    `db:"user_id"`
	Code       string    `db:"code"`
	IsVerified string    `db:"is_verified"`
	CreatedAt  time.Time `db:"created_at"`
	ExpiresAt  time.Time `db:"expires_at"`
}
