package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var testRouter *mux.Router

func TestMain(m *testing.M) {
	var session *mgo.Session
	var err error

	// Router and Session
	session, testRouter, err = setup(DBURL)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}

	database = session.DB(DBTESTNAME)

	// Inserting currency data into testing db
	setupTestingData()

	m.Run()

	database.DropDatabase()
}

func setupTestingData() {
	file, err := ioutil.ReadFile("./testdata/dummyCurrency.json")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
	}

	var rate CurrencyRate

	err = json.Unmarshal(file, &rate)
	if err != nil {
		fmt.Printf("Failed to parse json: %v\n", err)
	}

	// Inserts the data 7 times for the sake of testing the average handler
	for i := 0; i < 7; i++ {
		insertData("rates", rate)
	}

}
