package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	//postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/samanthreddys/myweb.com/hash"
	"github.com/samanthreddys/myweb.com/rand"
)

//User struct for user model
type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"not null; unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}


//userService struct
type userService struct {
	UserDB
}


// userGorm struct

type userGorm struct{
	db *gorm.DB
	hmac hash.HMAC

}

//UserValidator struct
type UserValidator struct{
	UserDB
}

var _UserDB=&userGorm{}
var (
	//ErrNotFound not found in database
	ErrNotFound = errors.New("models: resource not found")
	//ErrInvalidID invaild id passed
	ErrInvalidID = errors.New("models: Invalid id passed")
	// ErrInvaildPassword Invalid password provided
	ErrInvaildPassword = errors.New("models: Incorrect Password provided")
	//ErrInvaildEmail Invalid email provided
	//ErrInvaildEmail = errors.New("models: Incorrect Email provided")
)

const (
	userPasswordPepper = "mysecretstring"
	hmacsecretkey      = "secrethmackey"
)
//UserDB is a interface
type UserDB interface {
	// methods to query single users by id, email and remember token
	ByID(id uint)(*User,error)
	LookByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(u *User) error
	Update(u *User) error
	Delete(id uint) error

	//Used to close a DB connection
	Close() error

	//Migration Helpers
	DestructiveReset() error
	AutoMigrate() error

}
// UserService is as set of methods used to manipulate and work with user model
type UserService interface {
	//Authenticate will verify if the provided user and password are correct.
	Authenticate(email, password string) (*User, error)
	UserDB
}

func (uv *UserValidator) ByID(id uint)(*User, error){
	// Validate ID , if id<0 return invalid id
	if id<=0{
			return nil, errors.New("Invalid id")

	}
	return uv.UserDB.ByID(id)
}

func newUserGorm(connectioninfo string)(*userGorm, error){
	db, err := gorm.Open("postgres", connectioninfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacsecretkey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
	}

//NewUserService to connect to db
func NewUserService(connectioninfo string) (UserService, error) {
	ug,err:=newUserGorm(connectioninfo)
	if err!=nil{
		return nil,err
	}
	return &userService{
		UserDB:&UserValidator{

			UserDB:ug,
		},
	}, nil

}

//Close the db connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset is used to drop table if it exists
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		panic(err)
	}
	return ug.AutoMigrate()
}

//AutoMigrate will attempt to automatically migrate users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		panic(err)
	}

	return nil
}

//Create user in user
func (ug *userGorm) Create(u *User) error {
	pwBytes := []byte(u.Password + userPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedBytes)
	u.Password = ""
	if u.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		u.Remember = token

	}
	u.RememberHash = ug.hmac.Hash(u.Remember)

	return ug.db.Create(u).Error
}

//Update user in user
func (ug *userGorm) Update(u *User) error {
	if u.Remember != "" {
		u.RememberHash = ug.hmac.Hash(u.Remember)
	}

	return ug.db.Save(u).Error
}

//Delete user in user
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	u := User{
		Model: gorm.Model{
			ID: id,
		},
	}
	return ug.db.Delete(&u).Error

}

//ByID will look up by id provided
// 1- user, nil
// 2- nil,error not found
// 3- nil, other error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id=?", id)
	err := first(db, &user)
	return &user, err

}

//ByRemember looks up user using a remember token
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	hashedToken := ug.hmac.Hash(token)
	db := ug.db.Where("remember_hash=?", hashedToken)
	err := first(db, &user)
	return &user, err

}

//LookByEmail look user by email
func (ug *userGorm) LookByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email=?", email)
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

// Authenticate method is used to authenticate user login It returns user details or error incase of failure.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.LookByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPasswordPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvaildPassword
		default:
			return nil, err

		}

	}
	return foundUser, nil
}
