package models

type Category struct {
	BaseModel

	Name string `gorm:"size:100;not null"`

	Description *string

	Courses []Course `gorm:"foreignKey:CategoryID"`
}
