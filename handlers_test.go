package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var testId string

//------------------------------------------------------------------------------
func TestLatest(t *testing.T) {
	request := Ticket{
		Base:   "EUR",
		Target: "NOK",
	}

	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(request)

	req, err := http.NewRequest("POST", "/latest", data)
	if err != nil {
		t.Fatal("Failed creating POST request")
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}

	//testing value against referance value
	file, err := ioutil.ReadFile("./testdata/dummyCurrency.json")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
	}

	var rate CurrencyRate
	err = json.Unmarshal(file, &rate)
	if err != nil {
		fmt.Printf("Failed to parse json: %v\n", err)
	}

	var result float64
	json.NewDecoder(recorder.Body).Decode(&result)

	if result != rate.Rates[request.Target] {
		t.Fatal("Average should be ", rate.Rates[request.Target], "   - returned ", result)
	}
}

//------------------------------------------------------------------------------
func TestAverage(t *testing.T) {
	request := Ticket{
		Base:   "EUR",
		Target: "NOK",
	}

	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(request)

	req, err := http.NewRequest("POST", "/average", data)
	if err != nil {
		t.Fatal("Failed creating POST request")
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}

	//testing value against referance value
	file, err := ioutil.ReadFile("./testdata/dummyCurrency.json")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
	}

	var rate CurrencyRate
	err = json.Unmarshal(file, &rate)
	if err != nil {
		fmt.Printf("Failed to parse json: %v\n", err)
	}

	var result float64
	json.NewDecoder(recorder.Body).Decode(&result)

	if result != rate.Rates[request.Target] {
		t.Fatal("Average should be ", rate.Rates[request.Target], "   - returned ", result)
	}
}

//------------------------------------------------------------------------------
func TestNewHook(t *testing.T) {
	request := Ticket{
		URL:        "http://remoteUrl:8080/randomWebhookPath",
		Base:       "EUR",
		Target:     "NOK",
		MinTrigger: 7,
		MaxTrigger: 11,
	}

	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(request)

	req, err := http.NewRequest("POST", "/", data)
	if err != nil {
		t.Fatal("Failed creating POST request")
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}

	result := (recorder.Body).String()

	// Is a valid hex id returned?
	if !bson.IsObjectIdHex(result) {
		t.Fatal("Server conversion error. Returned id is not valid ObjectIdHex: ", result)
	}

	testId = result
}

//------------------------------------------------------------------------------
func TestAccessHook(t *testing.T) {
	req, err := http.NewRequest("GET", "/"+testId, nil)
	if err != nil {
		t.Fatal("Failed creating GET request")
	}

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatal("Server error: Returned ", recorder.Code, " instead of ", http.StatusOK)
	}

	var result Ticket
	json.NewDecoder(recorder.Body).Decode(&result)

	var ticket Ticket
	collection := database.C("tickets")
	collection.FindId(bson.ObjectIdHex(testId)).One(&ticket)

	if err != nil {
		fmt.Printf("Failed to parse json: %v\n", err)
	}

	if ticket != result {
		t.Fatal("Server could not find a entry with _id: ", testId)
	}

}
