package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//------------------------------------------------------------------------------
func updateAndInvoke() {
	updateRates()
	invokeClients()
}

//------------------------------------------------------------------------------
func updateRates() {
	var data CurrencyRate

	err := json.Unmarshal(fetchRates(), &data)

	if err != nil {
		log.Fatal(err)
	} else {
		insertData("rates", data)
	}
}

//------------------------------------------------------------------------------
func fetchRates() []byte {
	res, err := http.Get(RATEURL)

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return body
}

//-----------------------------------------------------------------------------
func invokeClients() {
	var tickets []Ticket
	var rate CurrencyRate

	collection := session.DB("CurrencyDB").C("tickets")
	collection.Find(nil).All(&tickets)

	collection = session.DB("CurrencyDB").C("rates")
	collection.Find(nil).Sort("-_id").One(&rate)

	for _, ticket := range tickets {
		if outOfBounds(ticket, rate) {
			notifyClient(ticket, rate)
		}
	}

}

func forceInvokeClients() {
	var tickets []Ticket
	var rate CurrencyRate

	collection := session.DB("CurrencyDB").C("tickets")
	collection.Find(nil).All(&tickets)

	collection = session.DB("CurrencyDB").C("rates")
	collection.Find(nil).Sort("-_id").One(&rate)

	for _, ticket := range tickets {
		notifyClient(ticket, rate)
	}
}

func outOfBounds(t Ticket, r CurrencyRate) bool {
	return (r.Rates[t.Target] > t.MaxTrigger) || (r.Rates[t.Target] < t.MinTrigger)
}

func notifyClient(t Ticket, r CurrencyRate) {
	invoke := InvokeData{
		Base:       t.Base,
		Target:     t.Target,
		Current:    r.Rates[t.Target],
		MinTrigger: t.MinTrigger,
		MaxTrigger: t.MaxTrigger,
	}

	response, _ := json.MarshalIndent(invoke, "", "   ")
	http.Post(t.URL, "application/x-www-form-urlencoded", bytes.NewBuffer(response))
}

//------------------------------------------------------------------------------
func insertData(collectionName string, data interface{}) {
	collection := session.DB("CurrencyDB").C(collectionName)

	err := collection.Insert(data)
	if err != nil {
		panic(err)
	}
}
