package views

import (
	"html/template"
)

// View struct
type View struct {
	Template *template.Template
}

//NewView function
func NewView(files ...string) *View {

	files = append(files, "views/layouts/footer.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
	}

}
