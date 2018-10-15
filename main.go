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
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "<h1>Welcome to Web Server Development</h1>")

}
