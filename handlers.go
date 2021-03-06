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
	defer r.Body.Close()

	var ticket Ticket

	err := decoder.Decode(&ticket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

	} else {
		ticket.ID = bson.NewObjectId()

		insertData("tickets", ticket)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, ticket.ID.Hex())
	}
}

//------------------------------------------------------------------------------
func handleAccessHook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	var ticket Ticket

	if !bson.IsObjectIdHex(vars["id"]) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid ID")
	} else {
		collection := database.C("tickets")
		collection.FindId(bson.ObjectIdHex(vars["id"])).One(&ticket)

		response, err := json.MarshalIndent(ticket, "", "   ")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, string(response))
		}
	}
}

//------------------------------------------------------------------------------
func handleDeleteHook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	collection := database.C("tickets")
	collection.RemoveId(bson.ObjectIdHex(vars["id"]))
}

func handleLatest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rate CurrencyRate
	var ticket Ticket

	err := decoder.Decode(&ticket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	} else {
		collection := database.C("rates")

		collection.Find(nil).Sort("-_id").One(&rate)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, rate.Rates[ticket.Target])
	}
}

//------------------------------------------------------------------------------
func handleAverage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var rates []CurrencyRate
	var ticket Ticket
	var avg float64

	err := decoder.Decode(&ticket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	} else {
		collection := database.C("rates")
		collection.Find(nil).Sort("-_id").Limit(7).All(&rates)

		for _, rate := range rates {
			avg += rate.Rates[ticket.Target]
		}

		avg /= float64(len(rates))

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, avg)
	}
}

//------------------------------------------------------------------------------
func handleEvaluationTrigger(w http.ResponseWriter, r *http.Request) {
	forceInvokeClients()
}
