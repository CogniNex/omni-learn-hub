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
)

type UsersService struct {
	repo   repository.Users
	hasher hash.PasswordHasher
	otp    otp.Generator
	sms    sms.SMSClient
}

type Users interface {
	SignUp(ctx context.Context, input dto.UserSignUpInput) error
	GetOtp(ctx context.Context, request dto.UserGetOtpRequest) error
}

func NewUserService(repo repository.Users, hasher hash.PasswordHasher, otp otp.Generator, sms sms.SMSClient) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
		otp:    otp,
		sms:    sms,
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
	err = s.repo.Create(ctx, newUser)
	if err != nil {
		return fmt.Errorf("UserService - SignUp - s.repo.Create: %w", err)
	}

	return nil
}

func (s *UsersService) GetOtp(ctx context.Context, request dto.UserGetOtpRequest) error {
	otpCode := s.otp.RandomSecret()

	templates := s.sms.GetTemplates()

	messageWithOtp := fmt.Sprintf(templates.Registration, otpCode)
	err := s.sms.SendSMS(messageWithOtp, request.PhoneNumber)
	if err != nil {
		return err
	}

	return nil

}
