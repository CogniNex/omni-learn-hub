package entity

type BackendType string

const (
	Int      BackendType = "int"
	Varchar  BackendType = "varchar"
	Decimal  BackendType = "decimal"
	Text     BackendType = "text"
	Datetime BackendType = "datetime"
	Bool     BackendType = "bool"
)

type Attribute struct {
	AttributeID  int         `db:"attribute_id"`
	EntityTypeID int         `db:"entity_type_id"`
	PropertyName string      `db:"attribute_name"`
	BackendType  BackendType `db:"backend_type"`
	IsVisible    bool        `db:"is_visible"`
	IsRequired   bool        `db:"is_required"`
}
