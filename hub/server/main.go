package main

import (
	"encoding/json"
	"fmt" // formatting and printing values to the console.
	"log" // logging messages to the console.
	"math/rand"
	"net/http" // Used for build HTTP servers and clients.
)

type OutletData struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	FranchiseId int     `json:"franchise_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Region      string  `json:"region"`
	District    string  `json:"district"`
	State       string  `json:"state"`
}

type Data struct {
	Version     string     `json:"version"`
	DateTime    string     `json:"datetime"`
	AppEnv      string     `json:"app_env"`
	Env         string     `json:"env"`
	Timezone    string     `string:"timezone"`
	Commit      string     `json:"commit"`
	DelegateJob string     `json:"delegate_job"`
	Outlet      OutletData `json:"outlet"`
}

type HubHealth struct {
	STATUS string `json:"status"`
}

// Port we listen on.
const portNum string = ":3000"

// Handler functions.
func Station(w http.ResponseWriter, r *http.Request) {
	outletData := OutletData{
		Name:        "SH SEDUAN LAND DISTRICT SIBU",
		Id:          3,
		Code:        "1660",
		FranchiseId: 5,
		Latitude:    2.301458,
		Longitude:   111.880804,
		Region:      "BORNEO",
		District:    "SIBU",
		State:       "SARAWAK",
	}

	response := Data{
		Version:     "develop-20240903-0010975ce9",
		DateTime:    "2024-09-03T07:46:42+08:00",
		AppEnv:      "SHELL",
		Env:         "staging",
		Timezone:    "Asia/Kuala_Lumpur",
		Commit:      "0010975ce9",
		DelegateJob: "true",
		Outlet:      outletData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func randomBool() string {
	if rand.Intn(100) < 80 {
		return "up"
	} else {
		return "down"
	}
}

func SipStatus(w http.ResponseWriter, r *http.Request) {
	response := HubHealth{
		STATUS: randomBool(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CapillaryStatus(w http.ResponseWriter, r *http.Request) {
	response := HubHealth{
		STATUS: randomBool(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func GhlStatus(w http.ResponseWriter, r *http.Request) {
	response := HubHealth{
		STATUS: randomBool(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	log.Println("Starting our simple http server.")

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/", Station)
	http.HandleFunc("/internal/health/shell_sip", SipStatus)
	http.HandleFunc("/internal/health/shell_capillary", CapillaryStatus)
	http.HandleFunc("/internal/health/shell_ghl", GhlStatus)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")

	// Spinning up the server.
	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
