package controllers

import (
	"fmt"
	"net/http"

	"github.com/samanthreddys/myweb.com/views"
)

//Users struct
type Users struct {
	NewView *views.View
}

// SignUpForm struct to hold sign up form values
type SignUpForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//NewUsers function
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/signup.gohtml"),
	}
}

//New  GET/signup...
// New Method to GET signup This is used to render a form when user click on a signup form
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Is this called?")
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)

	}

	//	u.NewView.Render(w, nil)

}

//Create Used to create user
//POST SIGNUP
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignUpForm
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form)
	fmt.Fprintln(w, r.PostForm["email"])
	fmt.Fprintln(w, r.PostForm["password"])

}
