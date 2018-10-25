package views

import (
	"html/template"
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
