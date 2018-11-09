package models

import (
	"github.com/jinzhu/gorm"
)

//User struct for user model
type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string `gorm:"not null; unique_index"`
}
