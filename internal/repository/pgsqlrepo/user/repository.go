package user

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"log"
	"omni-learn-hub/internal/domain/entity"
	"omni-learn-hub/pkg/postgres"
)

type UsersRepo struct {
	db *postgres.Postgres
}

func NewUsersRepo(pg *postgres.Postgres) *UsersRepo {
	return &UsersRepo{pg}
}

func (r *UsersRepo) Create(ctx context.Context, user entity.User, userProfile entity.UserProfile, roleId int) error {
	id, _ := uuid.NewV4()

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
		Columns("user_id, phone_number, password_hash, password_salt").
		Values(id, user.PhoneNumber, user.PasswordHash, user.PasswordSalt).
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
		Values(id, roleId).
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
		Columns("user_id, first_name, last_name, entity_type_id").
		Values(id, userProfile.FirstName, userProfile.Lastname, roleId).
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

func (r *UsersRepo) IsExist(ctx context.Context, phoneNumber string) (bool, error) {
	sql, args, err := r.db.Builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"phone_number": phoneNumber}).
		ToSql()

	row := r.db.Pool.QueryRow(ctx, sql, args...)
	var user entity.User
	err = row.Scan(&user.UserID, &user.PhoneNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return true, nil

}
