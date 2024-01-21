package repository

import (
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/entity"
)

type Users interface {
	Create(ctx context.Context, user entity.User) error
}
