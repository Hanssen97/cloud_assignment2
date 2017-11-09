package main

import (
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var testSession *mgo.Session
var testRouter *mux.Router

func TestMain(m *testing.M) {
	// Router and Session
	setupTesting()

	// Inserting currency data into testing db
	setupTestingData()

	m.Run()
}

func setupTesting() {
	var err error
	testSession, testRouter, err = setup()
	if err != nil {
		panic(err)
	}
}

func setupTestingData() {

}
