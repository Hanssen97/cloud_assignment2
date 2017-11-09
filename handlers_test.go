package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAverage(t *testing.T) {
	hook := Ticket{
		Base:   "EUR",
		Target: "NOK",
	}

	data := new(bytes.Buffer)
	json.NewEncoder(data).Encode(hook)

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

	if result != rate.Rates[hook.Target] {
		t.Fatal("Average should be ", rate.Rates[hook.Target], "   - returned ", result)
	}
}
