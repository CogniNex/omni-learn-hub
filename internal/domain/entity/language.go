package entity

type Language struct {
	LanguageID   int    `db:"language_id"`
	LanguageCode string `db:"language_code"`
	LanguageName string `db:"language_name"`
}
