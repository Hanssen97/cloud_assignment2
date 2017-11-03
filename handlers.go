package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

//------------------------------------------------------------------------------
func handleNewHook(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ticket Ticket

	err := decoder.Decode(&ticket)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	ticket.ID = bson.NewObjectId()

	insertData("tickets", ticket)

	fmt.Fprint(w, ticket.ID.Hex())
}

//------------------------------------------------------------------------------
func handleAccessHook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	var ticket Ticket

	collection := session.DB("CurrencyDB").C("tickets")
	collection.FindId(bson.ObjectIdHex(vars["id"])).One(&ticket)

	response, err := json.MarshalIndent(ticket, "", "   ")

	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, string(response))
	}
}

//------------------------------------------------------------------------------
func handleDeleteHook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	collection := session.DB("CurrencyDB").C("tickets")
	collection.RemoveId(bson.ObjectIdHex(vars["id"]))
}

func handleLatest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var rate CurrencyRate
	var ticket Ticket

	err := decoder.Decode(&ticket)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	collection := session.DB("CurrencyDB").C("rates")

	collection.Find(nil).Sort("-_id").One(&rate)

	fmt.Fprint(w, rate.Rates[ticket.Target])
}

//------------------------------------------------------------------------------
func handleAverage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var rates []CurrencyRate
	var ticket Ticket
	var avg float64

	err := decoder.Decode(&ticket)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()

	collection := session.DB("CurrencyDB").C("rates")
	collection.Find(nil).Sort("-_id").Limit(1).All(&rates)

	for _, rate := range rates {
		avg += rate.Rates[ticket.Target]
	}

	avg /= float64(len(rates))

	fmt.Fprint(w, avg)

}

//------------------------------------------------------------------------------
func handleEvaluationTrigger(w http.ResponseWriter, r *http.Request) {
	forceInvokeClients()
}
