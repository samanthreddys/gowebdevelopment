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
	FirstName string
	LastName  string
	Email     string `gorm:"not null; unique_index"`
}

var (
	//ErrNotFound not found in database
	ErrNotFound = errors.New("models: resource not found")
	//ErrInvalidID invaild id passed
	ErrInvalidID = errors.New("models: Invalid id passed")
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
	db.LogMode(true)
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

//Delete user in user
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	u := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	return us.db.Delete(&u).Error

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

// Function to query db and return first record to dst , if nothing found in query return error not found
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	default:
		return err
	}

}
