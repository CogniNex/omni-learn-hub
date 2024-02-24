package user

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"log"
	"omni-learn-hub/internal/domain/entity"
	tokenResponse "omni-learn-hub/internal/service/token/dto/response"
	userResponse "omni-learn-hub/internal/service/user/dto/response"
	"omni-learn-hub/pkg/postgres"
	"time"
)

type UsersRepo struct {
	db *postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) Create(ctx context.Context, user entity.User, userProfile entity.UserProfile, roleId int, token tokenResponse.TokenResponse) error {

	// Begin transaction
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - Begin transaction: %w", err)
	}
	defer func() {
		// Rollback transaction if there's an error or defer completes
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		if err := tx.Commit(ctx); err != nil {
			log.Printf("Error committing transaction: %v\n", err)
		}
	}()

	// Insert into users table
	sqlUser, argsUser, err := r.db.Builder.
		Insert("users").
		Columns("user_id, phone_number, password_hash, password_salt, refresh_token, refresh_expires_in").
		Values(user.UserID, user.PhoneNumber, user.PasswordHash, user.PasswordSalt, token.RefreshToken, time.Unix(token.RefreshTokenExpireTime, 0)).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder (users): %w", err)
	}

	_, err = tx.Exec(ctx, sqlUser, argsUser...)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - tx.Exec (users): %w", err)
	}

	// Insert into user_roles table
	sqlUserRole, argsUser, err := r.db.Builder.
		Insert("user_roles").
		Columns("user_id, role_id").
		Values(user.UserID, roleId).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder (user_role): %w", err)
	}

	_, err = tx.Exec(ctx, sqlUserRole, argsUser...)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - tx.Exec (user_role): %w", err)
	}

	// Insert into user_profile table
	sqlProfile, argsProfile, err := r.db.Builder.
		Insert("user_profiles").
		Columns("user_id, first_name, last_name, entity_type_id, is_active").
		Values(userProfile.UserID, userProfile.FirstName, userProfile.Lastname, roleId, true).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - Create - r.Builder (user_profile): %w", err)
	}

	_, err = tx.Exec(ctx, sqlProfile, argsProfile...)
	if err != nil {
		return fmt.Errorf("UserRepo - Create - tx.Exec (user_profile): %w", err)
	}

	return nil

}

func (r *UsersRepo) Get(ctx context.Context, phoneNumber string) (userResponse.UserDetails, error) {
	sql, args, err := r.db.Builder.
		Select("users.user_id", "users.phone_number", "users.password_hash", "users.password_salt",
			"user_profiles.first_name", "user_profiles.last_name", "user_profiles.is_active").
		From("users").
		Join("user_profiles ON users.user_id = user_profiles.user_id").
		Where(squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()

	row := r.db.Pool.QueryRow(ctx, sql, args...)
	var userDetails userResponse.UserDetails
	err = row.Scan(&userDetails.UserID, &userDetails.PhoneNumber, &userDetails.PasswordHash,
		&userDetails.PasswordSalt, &userDetails.FirstName, &userDetails.Lastname, &userDetails.IsActive)

	if errors.Is(err, pgx.ErrNoRows) {
		return userResponse.UserDetails{}, nil
	}
	return userDetails, nil

}

func (r *UsersRepo) RefreshToken(ctx context.Context, phoneNumber string, token tokenResponse.TokenResponse) error {
	// Build the SQL query using db.Builder
	sql, args, err := r.db.Builder.
		Update("users").
		SetMap(map[string]interface{}{
			"refresh_token":      token.RefreshToken,
			"refresh_expires_in": time.Unix(token.RefreshTokenExpireTime, 0),
		}).
		Where(
			squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()

	if err != nil {
		return fmt.Errorf("UserRepo - RefreshToken - o.Builder: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - RefreshToken - o.Pool.Exec: %w", err)
	}

	return nil
}
