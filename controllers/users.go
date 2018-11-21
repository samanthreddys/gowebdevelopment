package controllers

import (
	"fmt"
	"net/http"

	"github.com/samanthreddys/myweb.com/rand"

	"github.com/samanthreddys/myweb.com/models"

	"github.com/samanthreddys/myweb.com/views"
)

//Users struct
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

// SignUpForm struct to hold sign up form values
type SignUpForm struct {
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
	Email     string `schema:"email"`
	Password  string `schema:"password"`
}

// SignInForm struct to hold sign in form values
type SignInForm struct {
	Email    string `schema:"email"`
	Password string `schema:"Password"`
}

//NewUsers function
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/signup"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

//Create Used to create user
//POST SIGNUP
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignUpForm

	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Password:  form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := u.SetCookie(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)

}

//Login POST Login ...
//Login function is used to verify the user login and succesfully login user if the valid authentication details are provided
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := SignInForm{}

	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address")
		case models.ErrInvaildPassword:
			fmt.Fprintln(w, "Invalid password provided")
		case nil:
			fmt.Fprintln(w, user)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = u.SetCookie(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

//CookieTest is a test function used to view the cookie saved
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", cookie.Value)
}

//SetCookie is used to set cookie for a web interaction
func (u *Users) SetCookie(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}

	}
	cookie := http.Cookie{
		Name:  "remember_token",
		Value: user.Remember,
	}
	http.SetCookie(w, &cookie)
	return nil

}
