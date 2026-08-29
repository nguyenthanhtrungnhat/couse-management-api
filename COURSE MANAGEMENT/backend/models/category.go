package models

import "github.com/google/uuid"

type Category struct {
	BaseModel

	Name        string `json:"name"`
	Description string `json:"description"`

	Courses []Course `json:"courses,omitempty"`
}

// TableName returns the database table name.
func (Category) TableName() string {
	return "categories"
}

// NewCategory creates a new Category model.
func NewCategory(
	name string,
	description string,
) *Category {
	return &Category{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		Name:        name,
		Description: description,
	}
}
