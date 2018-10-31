package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samanthreddys/myweb.com/controllers"
)

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>We could not find the page you requested for:(</h1><p>Please email us if you keep being sent to invalid page.</p>")
}

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(notfound)

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
