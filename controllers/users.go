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
	fmt.Fprintln(w, "This is creating a user account")
}
