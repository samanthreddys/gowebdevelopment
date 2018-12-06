package models

import (
	"github.com/pkg/errors"
	"github.com/samanthreddys/myweb.com/rand"
	"golang.org/x/crypto/bcrypt"
	"strings"

	"github.com/jinzhu/gorm"
	//postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/samanthreddys/myweb.com/hash"
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


}

//UserValidator struct
type UserValidator struct{
	UserDB
	hmac hash.HMAC
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


func (uv *UserValidator) ByID(id uint)(*User, error){
	// Validate ID , if id<0 return invalid id
	if id<=0{
			return nil, errors.New("Invalid id")

	}
	return uv.UserDB.ByID(id)
}



func (uv *UserValidator) bycryptPassword(u *User)error{
	if u.Password == ""{
		return nil
	}
	pwBytes := []byte(u.Password + userPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedBytes)
	u.Password = ""
	return nil
}

func (uv *UserValidator) hmacRemember(u *User)error{
	if u.Remember==""{
		return nil
	}
	u.RememberHash=uv.hmac.Hash(u.Remember)
	return nil

}

func (uv *UserValidator) defaultRemember(u *User) error{
	if u.Remember != "" {
		return nil
	}

		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		u.Remember = token
		return nil

}
func (uv *UserValidator) idGreaterThan(n uint) userValidatorFunc{
	return userValidatorFunc(func(u *User)error{
		if u.ID<=n{
			return ErrInvaildPassword
		}
		return nil
	})
}

func (uv *UserValidator) normalizeEmail(u *User)error{
	u.Email= strings.ToLower(u.Email)
	u.Email= strings.TrimSpace(u.Email)
	return nil
}



type userValidatorFunc func(*User) error

//runuserValidatiorFunc
func runuserValidatorFunc(u *User,fns ...userValidatorFunc) error {
	for _,fn := range fns{
		if err:=fn(u);err!=nil{
			return err
		}
	}
	return nil
}
//Create func for user validatior
func (uv *UserValidator) Create(u *User) error {

	if err:= runuserValidatorFunc(u,uv.bycryptPassword,uv.defaultRemember,uv.hmacRemember,uv.normalizeEmail);err!=nil{
		return err
	}


	return uv.UserDB.Create(u)
}

//Update user in user
func (uv *UserValidator) Update(u *User) error {

	if err:= runuserValidatorFunc(u,uv.bycryptPassword,uv.defaultRemember,uv.hmacRemember,uv.normalizeEmail);err!=nil{
		return err
	}

	return uv.UserDB.Update(u)
}

//LookByEmail will normailze email address to make it to lower case and remove spaces
func (uv *UserValidator) LookByEmail(email string)( *User,error){
	u:=User{
		Email:email,
	}
	if err:=runuserValidatorFunc(&u,uv.normalizeEmail);err!=nil{
		return nil,err
	}
	return uv.UserDB.LookByEmail(u.Email)
}

//Delete user in user
func (uv *UserValidator) Delete(id uint) error{
	var u User
	u.ID=id

	if err:= runuserValidatorFunc(&u,uv.idGreaterThan(0));err!=nil{
		return nil
	}

	return uv.UserDB.Delete(id)

}

func (uv *UserValidator) ByRemember(token string)(*User,error){
	user:=User{
		Remember:token,
	}
	if err:= runuserValidatorFunc(&user,uv.hmacRemember);err!=nil{
		return nil, err
	}
	return	uv.UserDB.ByRemember(user.RememberHash)
}

func newUserGorm(connectioninfo string)(*userGorm, error){
	db, err := gorm.Open("postgres", connectioninfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &userGorm{
		db:   db,

	}, nil
	}

//NewUserService to connect to db
func NewUserService(connectioninfo string) (UserService, error) {
	ug,err:=newUserGorm(connectioninfo)
	if err!=nil{
		return nil,err
	}
	hmac := hash.NewHMAC(hmacsecretkey)
	return &userService{
		UserDB:&UserValidator{

			UserDB:ug,
			hmac:hmac,
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

	return ug.db.Create(u).Error
}

//Update user in user
func (ug *userGorm) Update(u *User) error {
	return ug.db.Save(u).Error
}

//Delete user in user
func (ug *userGorm) Delete(id uint) error {

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
func (ug *userGorm) ByRemember(rememberhash string) (*User, error) {
	var user User

	db := ug.db.Where("remember_hash=?", rememberhash)
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
