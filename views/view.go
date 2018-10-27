package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// View struct
type View struct {
	Template *template.Template
	Layout   string
}

var (
	layoutDir         = "views/layouts/"
	templateExtension = ".gohtml"
)

// Render is used to render view
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//NewView function
func NewView(layout string, files ...string) *View {

	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}

}

//layout files returns a slice of strings representing files
func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExtension)
	if err != nil {
		panic(err)
	}
	return files

}
