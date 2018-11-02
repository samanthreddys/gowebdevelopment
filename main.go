package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samanthreddys/myweb.com/controllers"
)

/* func pagenotfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, "views/static/pagenotfound.gohtml")

} */

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/pagenotfound", staticC.PageNotFound).Methods("GET")
	//r.NotFoundHandler = http.HandlerFunc(pagenotfound)
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
