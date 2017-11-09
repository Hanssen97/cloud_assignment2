package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/robfig/cron.v2"

	"github.com/gorilla/mux"
)

var (
	session *mgo.Session
)

//------------------------------------------------------------------------------
func main() {
	var router *mux.Router

	session, router = setup()

	updateRates()

	fmt.Println("App running on port " + PORT)

	http.ListenAndServe(":"+PORT, router)

	session.Close()
}

//------------------------------------------------------------------------------
func setup() (*mgo.Session, *mux.Router) {
	router := mux.NewRouter()
	session := setDataBase()
	setHandlers(router)
	setCrons()

	return session, router
}

//------------------------------------------------------------------------------
func setDataBase() *mgo.Session {
	session, err := mgo.Dial(DBURL)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	return session
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
