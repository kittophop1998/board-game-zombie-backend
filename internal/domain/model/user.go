package model

import "github.com/google/uuid"

type User struct {
	Model
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	FirstName string    `gorm:"column:first_name" json:"first_name"`
	LastName  string    `gorm:"column:last_name" json:"last_name"`
	Email     string    `gorm:"column:email;uniqueIndex" json:"email"`
}
