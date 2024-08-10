package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		return result, fmt.Errorf("Error performing GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("Error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("Error parsing JSON: %v", err)
	}
	return result, nil
}
