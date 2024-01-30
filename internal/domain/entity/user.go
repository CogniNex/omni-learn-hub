package entity

type User struct {
	UserID       string `db:"user_id"`
	PhoneNumber  string `db:"phone_number"`
	PasswordHash string `db:"password_hash"`
	PasswordSalt string `db:"password_salt"`
	RefreshToken string `db:"refresh_token"`
}
