package controllers

import (
	"github.com/samanthreddys/myweb.com/views"
)

//Static struct for  static controller
type Static struct {
	Home    *views.View
	Contact *views.View
}

//NewStatic function returns static views
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
	}

}
