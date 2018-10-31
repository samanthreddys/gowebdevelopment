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
	templateDir       = "views/"
	layoutDir         = "views/layouts/"
	templateExtension = ".gohtml"
)

// this function takes slice of strings representing filepath, prepends template directory to each string in slice
//eg: {home} would result {view/home} if templatedir =views/
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}

}

// this function takes slice of strings representing filepath for templates, appends template extension to each string in slice
//eg: {home} would result {home.gohtml} if templateext = .gohtml
func addTemplateExtFile(files []string) {
	for i, f := range files {
		files[i] = f + templateExtension
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}

}

// Render is used to render view
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//NewView function
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExtFile(files)
	//fmt.Println(files)

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
