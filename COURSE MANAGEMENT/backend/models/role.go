package models

type Role struct {
	BaseModel

	Name string `gorm:"size:50;uniqueIndex;not null"`

	Users []User `gorm:"foreignKey:RoleID"`
}