package entity

import "database/sql"

type EntityAttribute struct {
	EntityAttributeID int             `db:"entity_attribute_id"`
	EntityTypeID      int             `db:"entity_type_id"`
	EntityID          int             `db:"entity_id"`
	AttributeID       int             `db:"attribute_id"`
	AttrNum           int             `db:"attr_num"`
	AttrVarchar       sql.NullString  `db:"attr_varchar"`
	AttrBool          sql.NullBool    `db:"attr_bool"`
	AttrDatetime      sql.NullTime    `db:"attr_datetime"`
	AttrDecimal       sql.NullFloat64 `db:"attr_decimal"`
}
