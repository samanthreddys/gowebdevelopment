package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	//postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//User struct for user model
type User struct {
	gorm.Model
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `gorm:"not null; unique_index" json:"email,omitempty"`
}

var (
	//ErrNotFound not found in database
	ErrNotFound = errors.New("models: resource not found")
)

//UserService struct
type UserService struct {
	db *gorm.DB
}

//NewUserService to connect to db
func NewUserService(connectioninfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectioninfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(false)
	return &UserService{
		db: db,
	}, nil

}

//Close the db connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset is used to drop table if it exists
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

//Create user in user
func (us *UserService) Create(u *User) error {

	return us.db.Create(u).Error
}

//Update user in user
func (us *UserService) Update(u *User) error {

	return us.db.Save(u).Error
}

//ByID will look up by id provided
// 1- user, nil
// 2- nil,error not found
// 3- nil, other error
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id=?", id)
	err := first(db, &user)
	return &user, err

}

//LookByEmail look user by email
func (us *UserService) LookByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err

}

// Function to return a single record from the input request
func first(db *gorm.DB, user *User) error {
	err := db.First(user).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	default:
		return err
	}

}
