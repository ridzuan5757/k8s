package main

import (
	"encoding/json"
	"fmt" // formatting and printing values to the console.
	"log" // logging messages to the console.
	"math/rand"
	"net/http" // Used for build HTTP servers and clients.
)

type Data struct {
	SITE_NAME string `json:"site_name"`
	ROID      int    `json:"roid"`
	REGION    string `json:"region"`
	POSTCODE  int    `json:"postcode"`
	CITY      string `json:"city"`
	STATE     string `json:"state"`
}

type Hub struct {
	SIP       bool `json:"sip"`
	CAPILLARY bool `json:"capillary"`
	GHL       bool `json:"ghl"`
}

// Port we listen on.
const portNum string = ":8080"

// Handler functions.
func Station(w http.ResponseWriter, r *http.Request) {
	response := Data{
		SITE_NAME: "SH SEDUAN LAND DISTRICT SIBU",
		ROID:      1660,
		REGION:    "BORNEO",
		POSTCODE:  96000,
		CITY:      "SIBU",
		STATE:     "SARAWAK",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func randomBool() bool {
	return rand.Intn(100) < 80
}

func HubStatus(w http.ResponseWriter, r *http.Request) {
	response := Hub{
		SIP:       randomBool(),
		CAPILLARY: randomBool(),
		GHL:       randomBool(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	log.Println("Starting our simple http server.")

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/station", Station)
	http.HandleFunc("/hub", HubStatus)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")

	// Spinning up the server.
	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
