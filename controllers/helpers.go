package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

// ParseForm used to parse forms
func ParseForm(r *http.Request, destinationForm interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err

	}
	//var form SignUpForm
	decoder := schema.NewDecoder()
	if err := decoder.Decode(destinationForm, r.PostForm); err != nil {
		return err
	}
	return nil
}
