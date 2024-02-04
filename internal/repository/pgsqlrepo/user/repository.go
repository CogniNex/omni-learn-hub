package user

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/pkg/postgres"
)

type UsersRepo struct {
	db *postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) Create(ctx context.Context, user entity.User) error {
	id, _ := uuid.NewV4()
	sql, args, err := r.db.Builder.
		Insert("users").
		Columns("user_id, phone_number, password_hash, password_salt").
		Values(id, user.PhoneNumber, user.PasswordHash, user.PasswordSalt).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *UsersRepo) IsExist(ctx context.Context, phoneNumber string) (bool, error) {
	sql, args, err := r.db.Builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()

	row := r.db.Pool.QueryRow(ctx, sql, args...)
	var user entity.User
	err = row.Scan(&user.UserID, &user.PhoneNumber)

	// it doesn't work, but we should find the reason of it
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return true, nil

}
