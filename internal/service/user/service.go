package user

import (
	"fmt"
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/internal/repository"
	"omni-learn-hub/internal/service/user/dto"
	"omni-learn-hub/pkg/hash"
	"omni-learn-hub/pkg/otp"
	"omni-learn-hub/pkg/sms"
	"time"
)

type UsersService struct {
	usersRepo    repository.Users
	otpCodesRepo repository.OtpCodes
	hasher       hash.PasswordHasher
	otp          otp.Generator
	sms          sms.SMSClient
}

type Users interface {
	SignUp(ctx context.Context, input dto.UserSignUpInput) error
	GetOtp(ctx context.Context, request dto.UserGetOtpRequest) error
}

func NewUserService(usersRepo repository.Users, otpCodesRepo repository.OtpCodes, hasher hash.PasswordHasher, otp otp.Generator, sms sms.SMSClient) *UsersService {
	return &UsersService{
		usersRepo:    usersRepo,
		otpCodesRepo: otpCodesRepo,
		hasher:       hasher,
		otp:          otp,
		sms:          sms,
	}
}

func (s *UsersService) SignUp(ctx context.Context, input dto.UserSignUpInput) error {
	hashed_pwd, salt, err := s.hasher.HashPassword(input.Password)
	if err != nil {
		return err
	}
	newUser := entity.User{
		PhoneNumber:  input.PhoneNumber,
		PasswordHash: hashed_pwd,
		PasswordSalt: salt,
	}
	err = s.usersRepo.Create(ctx, newUser)
	if err != nil {
		return fmt.Errorf("UserService - SignUp - s.repoUsers.Create: %w", err)
	}

	return nil
}

func (s *UsersService) GetOtp(ctx context.Context, request dto.UserGetOtpRequest) error {

	isBlockedUser, err := s.isUserInBlackList(ctx, request.PhoneNumber)

	if err != nil {
		return fmt.Errorf("UserService - GetOtp - s.isUserInBlackList: %w", err)
	}
	if isBlockedUser {
		return fmt.Errorf("OTP generation is locked for this user: %w", err)
	}

	alreadyExistedValidOtp, err := s.otpCodesRepo.GetLastValidOtpByNumber(ctx, request.PhoneNumber)

	if err != nil {
		return fmt.Errorf("UserService - GetOtp - s.otpCodesRepo.GetLastValidOtpByNumber: %w", err)
	}

	if alreadyExistedValidOtp != (entity.OtpCode{}) && alreadyExistedValidOtp.GenerationAttempts >= 3 {
		err = s.otpCodesRepo.AddPhoneNumberToBlackList(ctx, request.PhoneNumber)
		if err != nil {
			return fmt.Errorf("UserService - GetOtp - s.otpCodesRepo.AddPhoneNumberToBlackList: %w", err)
		}
		return nil
	}

	err = s.generateOtpCode(ctx, request.PhoneNumber, alreadyExistedValidOtp)
	if err != nil {
		return fmt.Errorf("UserService - GetOtp - s.generateOtpCode: %w", err)
	}

	// logic for production
	//templates := s.sms.GetTemplates()

	//messageWithOtp := fmt.Sprintf(templates.Registration, otpCode)
	//err := s.sms.SendSMS(messageWithOtp, request.PhoneNumber)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (s *UsersService) isUserInBlackList(ctx context.Context, phoneNumber string) (bool, error) {
	blockedUser, err := s.otpCodesRepo.GetBlockedUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return false, fmt.Errorf("UserService - isUserBlocked - s.repoOtpCodes.GetBlockedUserByPhoneNumber: %w", err)
	}
	if (entity.OtpBlacklist{}) != blockedUser {
		if blockedUser.NextUnblockDate.Before(time.Now()) {
			err = s.otpCodesRepo.DeleteUserFromBlackList(ctx, phoneNumber)
			if err != nil {
				return false, fmt.Errorf("UserService - isUserBlocked - s.repoOtpCodes.DeleteUserFromBlackList: %w", err)
			}

		}
		return true, nil
	}
	return false, nil
}

func (s *UsersService) generateOtpCode(ctx context.Context, phoneNumber string,
	alreadyExistedOtpCode entity.OtpCode) error {
	otpCode := s.otp.RandomSecret()
	if alreadyExistedOtpCode != (entity.OtpCode{}) {
		err := s.otpCodesRepo.IncrementAttempts(ctx, alreadyExistedOtpCode.OtpID, otpCode)
		if err != nil {
			return fmt.Errorf("UserService - generateOtpCode - s.otpCodesRepo.IncrementAttempts: %w", err)

		}
		return nil
	}

	newOtpCode := entity.OtpCode{
		PhoneNumber: phoneNumber,
		Code:        otpCode,
	}

	err := s.otpCodesRepo.Add(ctx, newOtpCode)
	if err != nil {
		return fmt.Errorf("UserService - generateOtpCode - s.repoOtpCodes.Add: %w", err)
	}
	return nil
}
