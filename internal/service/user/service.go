package user

import (
	"fmt"
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/base"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/internal/repository"
	"omni-learn-hub/internal/service/user/dto/request"
	"omni-learn-hub/internal/service/user/dto/response"
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
	SignUp(ctx context.Context, input request.UserSignUpRequest) base.ApiValueResponse
	GetOtp(ctx context.Context, request request.UserGetOtpRequest) error
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

func (s *UsersService) SignUp(ctx context.Context, request request.UserSignUpRequest) base.ApiValueResponse {
	hashed_pwd, salt, err := s.hasher.HashPassword(request.Password)
	if err != nil {
		return base.NewApiValueResponseWithError("system_error")
	}

	if request.Password != request.PasswordVerification {
		return base.NewApiValueResponseWithError("UserService - SignUp - Passwords do not match")
	}

	isUserExist, err := s.usersRepo.IsExist(ctx, request.PhoneNumber)

	if err != nil {
		return base.NewApiValueResponseWithError("UserService - SignUp - s.usersRepo.IsExist")
	}

	if isUserExist {
		return base.NewApiValueResponseWithError("UserService - SignUp - user already exists")
	}

	isOtpCodeCorrect, err := s.isOtpCodeCorrect(ctx, request.PhoneNumber, request.OtpCode)

	if err != nil {
		return base.NewApiValueResponseWithError("UserService - SignUp - s.isOtpCodeCorrect")
	}

	if isOtpCodeCorrect == false {
		return base.NewApiValueResponseWithError("UserService - SignUp - otp code is not correct")
	}

	newUser := entity.User{
		PhoneNumber:  request.PhoneNumber,
		PasswordHash: hashed_pwd,
		PasswordSalt: salt,
	}
	err = s.usersRepo.Create(ctx, newUser)
	if err != nil {
		return base.NewApiValueResponseWithError("UserService - SignUp - s.repoUsers.Create")
	}

	return base.NewApiValueResponse(response.UserSignUpResponse{PhoneNumber: request.PhoneNumber})

}

func (s *UsersService) GetOtp(ctx context.Context, request request.UserGetOtpRequest) error {

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

func (s *UsersService) isOtpCodeCorrect(ctx context.Context, phoneNumber string, otpCode string) (bool, error) {

	validOtp, err := s.otpCodesRepo.VerifyOtp(ctx, phoneNumber, otpCode)
	if err != nil {
		return false, fmt.Errorf("UserService - isOtpCodeCorrect - s.repoOtpCodes.VerifyOtp: %w", err)
	}

	if validOtp == (entity.OtpCode{}) {
		return false, nil
	}

	err = s.otpCodesRepo.UpdateOtpVerification(ctx, validOtp.OtpID)

	if err != nil {
		return false, fmt.Errorf("UserService - isOtpCodeCorrect - s.repoOtpCodes.UpdateOtpVerification: %w", err)
	}

	return true, nil

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
				return true, fmt.Errorf("UserService - isUserBlocked - s.repoOtpCodes.DeleteUserFromBlackList: %w", err)
			}
			return false, nil
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
