package entity

import "time"

type OtpBlacklist struct {
	OtpBlackListId  int       `db:"otp_blacklist_id"`
	PhoneNumber     string    `db:"phone_number"`
	NextUnblockDate time.Time `db:"next_unblock_date"`
}
