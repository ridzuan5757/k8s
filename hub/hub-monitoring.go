package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	SIP       bool `json:"sip"`
	CAPILLARY bool `json:"capillary"`
	GHL       bool `json:"ghl"`
}

func getHub() (Response, error) {

	var result Response
	url := "http://localhost:8080/hub"

	resp, err := http.Get(url)
	if err != nil {
		return result, fmt.Errorf("Error performing GET request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("Error reading response body: ", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("Error parsing JSON: ", err)
	}
	return result, nil
}

func main() {
	response, err := getHub()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the parsed result
	fmt.Printf("SIP: %t\n", response.SIP)
	fmt.Printf("Capillary: %t\n", response.CAPILLARY)
	fmt.Printf("GHL: %t\n", response.GHL)
}
