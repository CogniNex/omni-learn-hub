package otp_code

import (
	sqlDb "database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/pkg/postgres"
	"time"
)

const _defaultEntityCap = 64

type OtpCodesRepo struct {
	db *postgres.Postgres
}

func NewOtpCodesRepo(pg *postgres.Postgres) *OtpCodesRepo {
	return &OtpCodesRepo{db: pg}
}

func (o *OtpCodesRepo) Add(ctx context.Context, otpCode entity.OtpCode) error {
	createdAt := time.Now()
	// Add 5 minuted to expiration time
	expiresAt := createdAt.Add(5 * time.Minute)
	sql, args, err := o.db.Builder.
		Insert("otp_codes").
		Columns("phone_number, code, is_verified, created_at, expires_at, generation_attempts").
		Values(otpCode.PhoneNumber, otpCode.Code, false, createdAt, expiresAt, 1).
		ToSql()

	if err != nil {
		return fmt.Errorf("OtpCodeRepo - Add - o.Builder: %w", err)
	}

	_, err = o.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OtpCodeRepo - Add - o.Pool.Exec: %w", err)
	}

	return nil
}

func (o *OtpCodesRepo) IsOtpBlackListed(ctx context.Context, phoneNumber string) (bool, error) {
	// Build the SQL query using db.Builder
	sql, args, err := o.db.Builder.
		Select("1").
		From("otp_blacklist").
		Where(squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("OtpCodeRepo - IsOtpBlackListed - o.Builder: %w", err)
	}

	_, err = o.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("OtpCodeRepo - IsOtpBlackListed - o.Pool.Exec: %w", err)
	}

	return true, nil
}

func (o *OtpCodesRepo) GetBlockedUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.OtpBlacklist, error) {
	sql, args, err := o.db.Builder.
		Select("*").
		From("otp_blacklist").
		Where(squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()
	if err != nil {
		return entity.OtpBlacklist{}, fmt.Errorf("OtpCodeRepo - GetBlockedUserByPhoneNumber - o.Builder: %w", err)
	}

	row := o.db.Pool.QueryRow(ctx, sql, args...)
	var otpBlackList entity.OtpBlacklist
	err = row.Scan(&otpBlackList.OtpBlackListId, &otpBlackList.PhoneNumber, &otpBlackList.NextUnblockDate)
	if errors.Is(err, sqlDb.ErrNoRows) {
		return entity.OtpBlacklist{}, fmt.Errorf("OtpCodeRepo - GetBlockedUserByPhoneNumber - r.Pool.Query: %w", err)
	}

	return otpBlackList, nil
}

func (o *OtpCodesRepo) DeleteUserFromBlackList(ctx context.Context, phoneNumber string) error {
	sql, args, err := o.db.Builder.
		Delete("otp_blacklist").
		Where(squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()

	if err != nil {
		return fmt.Errorf("OtpCodeRepo - Add - o.Builder: %w", err)
	}

	_, err = o.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OtpCodeRepo - Add - o.Pool.Exec: %w", err)
	}

	return nil
}

func (o *OtpCodesRepo) GetLastValidOtpByNumber(ctx context.Context, phoneNumber string) (entity.OtpCode, error) {

	currentTime := time.Now()
	sql, args, err := o.db.Builder.
		Select("*").
		From("otp_codes").
		Where(
			squirrel.And{
				squirrel.Eq{"phone_number": phoneNumber},
				squirrel.Gt{"expires_at": currentTime},
				squirrel.Eq{"is_verified": false},
			}).
		OrderBy("expires_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return entity.OtpCode{}, fmt.Errorf("OtpCodeRepo - GetLastValidOtpByNumber - o.Builder: %w", err)
	}

	row := o.db.Pool.QueryRow(ctx, sql, args...)
	var otpCode entity.OtpCode
	err = row.Scan(&otpCode.OtpID, &otpCode.Code, &otpCode.PhoneNumber, &otpCode.GenerationAttempts, &otpCode.IsVerified,
		&otpCode.CreatedAt, &otpCode.ExpiresAt)
	if errors.Is(err, sqlDb.ErrNoRows) {
		return entity.OtpCode{}, fmt.Errorf("OtpCodeRepo - GetBlockedUserByPhoneNumber - r.Pool.Query: %w", err)
	}

	return otpCode, nil
}

func (o *OtpCodesRepo) AddPhoneNumberToBlackList(ctx context.Context, phoneNumber string) error {
	currentTime := time.Now()
	// Add 5 minuted to expiration time
	nextUnblockDate := currentTime.Add(30 * time.Minute)
	otpBlackList := entity.OtpBlacklist{PhoneNumber: phoneNumber, NextUnblockDate: nextUnblockDate}
	sql, args, err := o.db.Builder.
		Insert("otp_blacklist").
		Columns("phone_number, next_unblock_date").
		Values(otpBlackList.PhoneNumber, otpBlackList.NextUnblockDate).
		ToSql()

	if err != nil {
		return fmt.Errorf("OtpCodeRepo - AddPhoneNumberToBlackList - o.Builder: %w", err)
	}

	_, err = o.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OtpCodeRepo - AddPhoneNumberToBlackList - o.Pool.Exec: %w", err)
	}

	return nil

}

func (o *OtpCodesRepo) IncrementAttempts(ctx context.Context, otpID int, otpCode string) error {

	// Build the SQL query using db.Builder
	sql, args, err := o.db.Builder.
		Update("otp_codes").
		SetMap(map[string]interface{}{
			"generation_attempts": squirrel.Expr("generation_attempts + 1"),
			"code":                otpCode,
		}).
		Where(
			squirrel.Eq{"otp_id": otpID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("OtpCodeRepo - IncrementAttempts - o.Builder: %w", err)
	}

	_, err = o.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OtpCodeRepo - IncrementAttempts - o.Pool.Exec: %w", err)
	}

	return nil
}

func (o *OtpCodesRepo) VerifyOtp(ctx context.Context, phoneNumber string, code string) (entity.OtpCode, error) {
	currentTime := time.Now()

	sql, args, err := o.db.Builder.
		Select("*").
		From("otp_codes").
		Where(
			squirrel.And{
				squirrel.Eq{"phone_number": phoneNumber},
				squirrel.Gt{"expires_at": currentTime},
				squirrel.Eq{"code": code},
			}).
		OrderBy("expires_at DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return entity.OtpCode{}, fmt.Errorf("OtpCodeRepo - VerifyOtp - o.Builder: %w", err)
	}

	row := o.db.Pool.QueryRow(ctx, sql, args...)
	var otpCode entity.OtpCode
	err = row.Scan(&otpCode.OtpID, &otpCode.Code, &otpCode.PhoneNumber, &otpCode.GenerationAttempts, &otpCode.IsVerified,
		&otpCode.CreatedAt, &otpCode.ExpiresAt)
	if errors.Is(err, sqlDb.ErrNoRows) {
		return entity.OtpCode{}, fmt.Errorf("OtpCodeRepo - VerifyOtp - r.Pool.Query: %w", err)
	}

	return otpCode, nil
}

func (o *OtpCodesRepo) UpdateOtpVerification(ctx context.Context, otpId int) error {
	// Build the SQL query using db.Builder
	sql, args, err := o.db.Builder.
		Update("otp_codes").
		SetMap(map[string]interface{}{
			"is_verified": true,
		}).
		Where(
			squirrel.Eq{"otp_id": otpId}).
		ToSql()

	if err != nil {
		return fmt.Errorf("OtpCodeRepo - UpdateOtpVerification - o.Builder: %w", err)
	}

	_, err = o.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OtpCodeRepo - UpdateOtpVerification - o.Pool.Exec: %w", err)
	}

	return nil
}
