package models

type Role struct {
	BaseModel

	Name string `json:"name"`

	Users []User `json:"users,omitempty"`
}