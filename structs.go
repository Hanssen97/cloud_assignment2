package main

import "gopkg.in/mgo.v2/bson"

// CurrencyRate holds raw currency rate data
type CurrencyRate struct {
	ID    bson.ObjectId      `bson:"_id,omitempty"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

// Ticket holds client request
type Ticket struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	URL        string        `json:"webhookURL"`
	Base       string        `json:"baseCurrency"`
	Target     string        `json:"targetCurrency"`
	MinTrigger float64       `json:"minTriggerValue"`
	MaxTrigger float64       `json:"maxTriggerValue"`
}

// InvokeData holds data to send client
type InvokeData struct {
	Base       string  `json:"baseCurrency"`
	Target     string  `json:"targetCurrency"`
	Current    float64 `json:"currentRate"`
	MinTrigger float64 `json:"minTriggerValue"`
	MaxTrigger float64 `json:"maxTriggerValue"`
}
