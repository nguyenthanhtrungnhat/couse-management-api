package models

import "github.com/google/uuid"

type Role struct {
	BaseModel

	Name string `json:"name"`

	Users []User `json:"users,omitempty"`
}

// TableName returns the database table name.
func (Role) TableName() string {
	return "roles"
}

// NewRole creates a new Role model.
func NewRole(name string) *Role {
	return &Role{
		BaseModel: BaseModel{
			ID: uuid.New(),
		},
		Name: name,
	}
}
