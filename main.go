package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/robfig/cron.v2"

	"github.com/gorilla/mux"
)

var (
	database *mgo.Database
)

//------------------------------------------------------------------------------
func main() {
	var router *mux.Router
	var session *mgo.Session
	var err error

	session, router, err = setup(DBURL)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	database = session.DB(DBNAME)

	updateRates()

	fmt.Println("App running on port " + PORT)

	http.ListenAndServe(":"+PORT, router)

	session.Close()
}

//------------------------------------------------------------------------------
func setup(url string) (*mgo.Session, *mux.Router, error) {
	router := mux.NewRouter()
	session, err := mgo.Dial(url)

	setHandlers(router)
	setCrons()

	return session, router, err
}

//------------------------------------------------------------------------------
func setCrons() {
	c := cron.New()
	c.AddFunc("@daily", updateAndInvoke)
	c.Start()
}

//------------------------------------------------------------------------------
func setHandlers(router *mux.Router) {
	router.HandleFunc("/", handleNewHook).Methods("POST")
	router.HandleFunc("/latest", handleLatest).Methods("POST")
	router.HandleFunc("/average", handleAverage).Methods("POST")
	router.HandleFunc("/evaluationtrigger", handleEvaluationTrigger).Methods("GET")
	router.HandleFunc("/{id}", handleAccessHook).Methods("GET")
	router.HandleFunc("/{id}", handleDeleteHook).Methods("DELETE")
}
