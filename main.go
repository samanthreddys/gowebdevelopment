package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

/* func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome to Web Server Development</h1>")

	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "To get in touch please email to <a href= \"mailto:test@test.com\"> test@test.com.</a>")

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>We could not find the page you requested for:(</h1><p>Please email us if you keep being sent to invalid page.</p>")
	}

} */

func faqs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "This ia a FAQ's Page. Work in Progress !")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to Web Server Development</h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch please email to <a href= \"mailto:test@test.com\"> test@test.com.</a>")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faqs", faqs)
	http.ListenAndServe(":3000", r)

}
