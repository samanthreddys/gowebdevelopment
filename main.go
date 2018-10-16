package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)

}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome to Web Server Development</h1>")

	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "To get in touch please email to <a href= \"mailto:test@test.com\"> test@test.com.</a>")

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>We could not find the page you requested for:(</h1><p>Please email us if you keep being sent to invalid page.</p>")
	}

}
