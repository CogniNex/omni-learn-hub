package entity

type Role struct {
	RoleID   int    `db:"role_id"`
	RoleName string `db:"role_name"`
}
