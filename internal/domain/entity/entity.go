package entity

type EntityType string

const (
	Teacher EntityType = "teacher"
	Student EntityType = "student"
	Parent  EntityType = "parent"
)

// Entities table
type Entity struct {
	EntityTypeID int        `db:"entity_type_id"`
	EntityType   EntityType `db:"entity_type"`
}
