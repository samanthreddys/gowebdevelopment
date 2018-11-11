package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samanthreddys/myweb.com/controllers"
	"github.com/samanthreddys/myweb.com/models"
)

/* func pagenotfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, "views/static/pagenotfound.gohtml")

} */

const (
	host   = "localhost"
	user   = "postgres"
	port   = 5432
	dbname = "clrpix_dev"
)

func main() {
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	us, err := models.NewUserService(psqlinfo)
	must(err)

	defer us.Close()
	us.AutoMigrate()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

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
