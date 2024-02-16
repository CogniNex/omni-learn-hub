package repository

import (
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/entity"
)

type Users interface {
	Create(ctx context.Context, user entity.User, userProfile entity.UserProfile, roleId int) error
	IsExist(ctx context.Context, phoneNumber string) (bool, error)
}

type OtpCodes interface {
	Add(ctx context.Context, otpCode entity.OtpCode) error
	IncrementAttempts(ctx context.Context, otpID int, otpCode string) error
	GetBlockedUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.OtpBlacklist, error)
	DeleteUserFromBlackList(ctx context.Context, phoneNumber string) error
	GetLastValidOtpByNumber(ctx context.Context, phoneNumber string) (entity.OtpCode, error)
	AddPhoneNumberToBlackList(ctx context.Context, phoneNumber string) error
	VerifyOtp(ctx context.Context, phoneNumber string, code string) (entity.OtpCode, error)
	UpdateOtpVerification(ctx context.Context, otpId int) error
}
