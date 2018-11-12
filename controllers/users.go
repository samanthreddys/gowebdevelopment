package controllers

import (
	"fmt"
	"net/http"

	"github.com/samanthreddys/myweb.com/models"

	"github.com/samanthreddys/myweb.com/views"
)

//Users struct
type Users struct {
	NewView *views.View
	us      *models.UserService
}

// SignUpForm struct to hold sign up form values
type SignUpForm struct {
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
	Email     string `schema:"email"`
	Password  string `schema:"password"`
}

//NewUsers function
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/signup"),
		us:      us,
	}
}

//New  GET/signup...
// New Method to GET signup This is used to render a form when user click on a signup form
func (u *Users) New(w http.ResponseWriter, r *http.Request) {

	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)

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

	fmt.Fprintln(w, user)

}
