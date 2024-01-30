package auth

import (
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/internal/service/common/models"
)

type TokenManager interface {
	CreateToken(user entity.User)
	ValidateToken(accessToken string)
	ValidateRefreshToken(model models.Token)
}
