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
	router  *mux.Router
)

//------------------------------------------------------------------------------
func main() {
	setup()

	updateRates()

	fmt.Println("App running on port " + PORT)

	http.ListenAndServe(":"+PORT, router)

	session.Close()
}

//------------------------------------------------------------------------------
func setup() {
	router = mux.NewRouter()

	setDataBase(DBURL)
	setCrons()
	setHandlers()
}

//------------------------------------------------------------------------------
func setDataBase(url string) {
	var err error
	session, err = mgo.Dial(url)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}

//------------------------------------------------------------------------------
func setCrons() {
	c := cron.New()
	c.AddFunc("@daily", updateAndInvoke)
	c.Start()
}

//------------------------------------------------------------------------------
func setHandlers() {
	router.HandleFunc("/", handleNewHook).Methods("POST")
	router.HandleFunc("/latest", handleLatest).Methods("POST")
	router.HandleFunc("/average", handleAverage).Methods("POST")
	router.HandleFunc("/evaluationtrigger", handleEvaluationTrigger).Methods("GET")
	router.HandleFunc("/{id}", handleAccessHook).Methods("GET")
	router.HandleFunc("/{id}", handleDeleteHook).Methods("DELETE")
}
