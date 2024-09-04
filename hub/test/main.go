package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// JSON object
	data := `{
		"version": "develop-20240903-0010975ce9",
		"datetime": "2024-09-03T07:46:42+08:00",
		"app_env": "SHELL",
		"env": "staging",
		"timezone": "Asia/Kuala_Lumpur",
		"commit": "0010975ce9",
		"delegate_job": "true",
		"outlet": {
			"id": 3,
			"name": "SH JLN SULTAN MAHMUD 2 KT x BANGI (POC)",
			"code": "1550",
			"franchise_id": 5,
			"latitude": "2.954362",
			"longitude": "101.758384",
			"region": "north",
			"state": "Selangor",
			"district": "Bandar Baru BANGI"
		}
	}`

	// Unmarshal JSON to a map
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(data), &jsonData)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Create or open .env file
	file, err := os.Create(".env")
	if err != nil {
		log.Fatalf("Error creating .env file: %v", err)
	}
	defer file.Close()

	// Process the JSON and write to .env file
	processJSON(jsonData, "", file)
}

// processJSON recursively processes JSON and writes to file using parent.child convention
func processJSON(data map[string]interface{}, parentKey string, file *os.File) {
	for key, value := range data {
		fullKey := key
		if parentKey != "" {
			fullKey = fmt.Sprintf("%s_%s", parentKey, key)
		}

		switch v := value.(type) {
		case map[string]interface{}:
			processJSON(v, fullKey, file) // Recurse for nested objects
		default:
			line := fmt.Sprintf("%s=%v\n", strings.ToUpper(fullKey), v)
			_, err := file.WriteString(line)
			if err != nil {
				log.Fatalf("Error writing to .env file: %v", err)
			}
		}
	}
}
