package user

import (
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/base"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/internal/repository"
	"omni-learn-hub/internal/service/token"
	"omni-learn-hub/internal/service/token/dto"
	"omni-learn-hub/internal/service/user/dto/request"
	"omni-learn-hub/internal/service/user/dto/response"
	"omni-learn-hub/pkg/hash"
	"omni-learn-hub/pkg/otp"
	"omni-learn-hub/pkg/sms"
	"omni-learn-hub/pkg/whatsapp"
	"time"
)

type UsersService struct {
	usersRepo       repository.Users
	otpCodesRepo    repository.OtpCodes
	tokenService    token.TokenService
	hasher          hash.PasswordHasher
	otp             otp.Generator
	sms             sms.SMSClient
	whatsappService whatsapp.WhatsappClient
}

type Users interface {
	SignUp(ctx context.Context, input request.UserSignUpRequest) base.ApiValueResponse
	GetOtp(ctx context.Context, request request.UserGetOtpRequest) base.ApiValueResponse
}

func NewUserService(usersRepo repository.Users, otpCodesRepo repository.OtpCodes, hasher hash.PasswordHasher,
	otp otp.Generator, sms sms.SMSClient, tokenService token.TokenService, whatsappService whatsapp.WhatsappClient) *UsersService {
	return &UsersService{
		usersRepo:       usersRepo,
		otpCodesRepo:    otpCodesRepo,
		hasher:          hasher,
		otp:             otp,
		sms:             sms,
		tokenService:    tokenService,
		whatsappService: whatsappService,
	}
}

func (s *UsersService) SignUp(ctx context.Context, request request.UserSignUpRequest) base.ApiValueResponse {

	hashed_pwd, salt, err := s.hasher.HashPassword(request.Password)
	if err != nil {
		return base.NewApiValueResponseWithError("system_error")
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

	id, _ := uuid.NewV4()

	newUser := entity.User{
		UserID:       id.String(),
		PhoneNumber:  request.PhoneNumber,
		PasswordHash: hashed_pwd,
		PasswordSalt: salt,
	}

	newUserProfile := entity.UserProfile{
		UserID:    id.String(),
		FirstName: request.FirstName,
		Lastname:  request.LastName,
	}

	tokenDto := dto.TokenDto{
		UserId:      id,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
	}
	tokenDto.Roles = append(tokenDto.Roles, "admin")

	tokenResponse, err := s.tokenService.GenerateToken(&tokenDto)

	if err != nil {
		return base.NewApiValueResponseWithError("TokenService - SignUp - s.tokenService.GenerateToken")
	}

	err = s.usersRepo.Create(ctx, newUser, newUserProfile, request.RoleId, *tokenResponse)
	if err != nil {
		return base.NewApiValueResponseWithError("UserService - SignUp - s.repoUsers.Create")
	}

	return base.NewApiValueResponse(tokenResponse)

}

func (s *UsersService) GetOtp(ctx context.Context, request request.UserGetOtpRequest) base.ApiValueResponse {

	isBlockedUser, err := s.isUserInBlackList(ctx, request.PhoneNumber)

	if err != nil {
		return base.NewApiValueResponseWithError("UserService - GetOtp - s.isUserInBlackList")
	}
	if isBlockedUser {
		return base.NewApiValueResponseWithError("OTP generation is locked for this user")
	}

	alreadyExistedValidOtp, err := s.otpCodesRepo.GetLastValidOtpByNumber(ctx, request.PhoneNumber)

	if err != nil {
		return base.NewApiValueResponseWithError("UserService - GetOtp - s.otpCodesRepo.GetLastValidOtpByNumber")
	}

	if alreadyExistedValidOtp != (entity.OtpCode{}) && alreadyExistedValidOtp.GenerationAttempts >= 3 {
		err = s.otpCodesRepo.AddPhoneNumberToBlackList(ctx, request.PhoneNumber)
		if err != nil {
			return base.NewApiValueResponseWithError("UserService - GetOtp - s.otpCodesRepo.AddPhoneNumberToBlackList")
		}
		return base.NewApiValueResponseWithError("You exceeded your current otp generation attempts")
	}

	otpCode, err := s.generateOtpCode(ctx, request.PhoneNumber, alreadyExistedValidOtp)
	if err != nil {
		return base.NewApiValueResponseWithError("UserService - GetOtp - s.generateOtpCode")
	}

	// logic for production
	templates := s.sms.GetTemplates()

	messageWithOtp := fmt.Sprintf(templates.Registration, otpCode)
	err = s.whatsappService.SendMessage(messageWithOtp, request.PhoneNumber)
	if err != nil {
		return base.NewApiValueResponseWithError("UserService - SendMessage - s.whatsappService.SendMessage")
	}

	return base.NewApiValueResponse(response.GetOtpResponse{PhoneNumber: request.PhoneNumber})
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
	alreadyExistedOtpCode entity.OtpCode) (string, error) {
	otpCode := s.otp.RandomSecret()
	if alreadyExistedOtpCode != (entity.OtpCode{}) {
		err := s.otpCodesRepo.IncrementAttempts(ctx, alreadyExistedOtpCode.OtpID, otpCode)
		if err != nil {
			return otpCode, fmt.Errorf("UserService - generateOtpCode - s.otpCodesRepo.IncrementAttempts: %w", err)

		}
		return otpCode, nil
	}

	newOtpCode := entity.OtpCode{
		PhoneNumber: phoneNumber,
		Code:        otpCode,
	}

	err := s.otpCodesRepo.Add(ctx, newOtpCode)
	if err != nil {
		return otpCode, fmt.Errorf("UserService - generateOtpCode - s.repoOtpCodes.Add: %w", err)
	}
	return otpCode, nil
}
