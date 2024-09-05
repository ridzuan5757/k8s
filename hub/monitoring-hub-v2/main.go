package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

func validateURL(endpoint string) error {
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return fmt.Errorf("Invalid url %v: %v", endpoint, err)
	}
	fmt.Printf("Using endpoint %v\n", endpoint)
	return nil
}

func fetchSiteData(endpoint string) (*Data, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error requesting site information: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error non-OK HTTP status: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	var outletData Data
	if err = json.Unmarshal(body, &outletData); err != nil {
		return nil, fmt.Errorf("Error parsing JSON: %v", err)
	}
	return &outletData, nil
}

func initContainers() error {
	cmd := exec.Command("docker-compose", "up", "-d")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error executing the command %v", err)
	}
	fmt.Printf("Command output:\n\r%v\n", string(output))
	return nil
}

func writeEnv(file *os.File, key, value interface{}) {
	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	case int, int64, float64, float32:
		valueStr = fmt.Sprintf("%v", v)
	default:
		log.Fatalf("Unsupported value type for key %s: %T", key, value)
	}
	_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, valueStr))
	if err != nil {
		log.Fatalf("Error writing to .env file: %v", err)
	}
}

func main() {
	endpoint := flag.String("endpoint", "https://default-endpoint", "The API endpoint to connect to")
	flag.Parse()

	if err := validateURL(*endpoint); err != nil {
		log.Fatalf("Invalid URL: %v\nExiting...\n", err)
	}

	outletData, err := fetchSiteData(*endpoint)
	if err != nil {
		log.Fatalf("Error fetching site data: %v\nExiting...\n", err)
	}
	fmt.Printf("Fetched data: %+v\n", outletData)

	file, err := os.Create(".env")
	if err != nil {
		log.Fatalf("Error creating .env file: %v", err)
	}
	defer file.Close()

	writeEnv(file, "OUTLET_NAME", outletData.Outlet.Name)
	writeEnv(file, "OUTLET_ID", outletData.Outlet.Id)
	writeEnv(file, "OUTLET_REGION", outletData.Outlet.Region)
	writeEnv(file, "OUTLET_DISTRICT", outletData.Outlet.District)
	writeEnv(file, "OUTLET_STATE", outletData.Outlet.State)
	writeEnv(file, "OUTLET_LATITUDE", outletData.Outlet.Latitude)
	writeEnv(file, "OUTLET_LONGITUDE", outletData.Outlet.Longitude)
	writeEnv(file, "ENVIRONMENT", outletData.Env)

	if err := initContainers(); err != nil {
		log.Fatalf("Error executing Docker Compose: %v", err)
	}
}
