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
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}

}
