package user

import (
	"fmt"
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/internal/repository"
	"omni-learn-hub/internal/service/user/dto"
	"omni-learn-hub/pkg/hash"
)

type UsersService struct {
	repo   repository.Users
	hasher hash.PasswordHasher
}

type Users interface {
	SignUp(ctx context.Context, input dto.UserSignUpInput) error
}

func NewUserService(repo repository.Users, hasher hash.PasswordHasher) *UsersService {
	return &UsersService{
		repo:   repo,
		hasher: hasher,
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
		return fmt.Errorf("TranslationUseCase - Translate - s.repo.Store: %w", err)
	}

	return nil
}
