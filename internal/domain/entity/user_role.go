package entity

type UserRole struct {
	UserID string `db:"user_id"`
	RoleID int    `db:"role_id"`
}
