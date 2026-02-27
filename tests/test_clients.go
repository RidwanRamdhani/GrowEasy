package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"GrowEasy/services/integration"
)

func main() {
	// Create output file
	file, err := os.Create("tests/output.txt")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	// Sample coordinates (Jakarta, Indonesia)
	lat := -6.2088
	lng := 106.8456

	// Test Weather Client
	fmt.Println("=== Testing Weather Client ===")
	file.WriteString("=== Testing Weather Client ===\n")

	weatherClient := integration.NewWeatherClient()
	weatherData, err := weatherClient.GetWeather(lat, lng)
	if err != nil {
		errMsg := fmt.Sprintf("Weather error: %v\n", err)
		fmt.Print(errMsg)
		file.WriteString(errMsg)
	} else {
		weatherJSON, _ := json.MarshalIndent(weatherData, "", "  ")
		fmt.Println(string(weatherJSON))
		file.WriteString(string(weatherJSON))
		file.WriteString("\n")
	}

	// Test Soil Client
	fmt.Println("\n=== Testing Soil Client ===")
	file.WriteString("\n=== Testing Soil Client ===\n")

	soilClient := integration.NewSoilClient()
	soilData, err := soilClient.GetSoilData(lat, lng)
	if err != nil {
		errMsg := fmt.Sprintf("Soil error: %v\n", err)
		fmt.Print(errMsg)
		file.WriteString(errMsg)
	} else {
		soilJSON, _ := json.MarshalIndent(soilData, "", "  ")
		fmt.Println(string(soilJSON))
		file.WriteString(string(soilJSON))
		file.WriteString("\n")
	}

	fmt.Println("\nHasil lengkap sudah disimpan di tests/output.txt")
	file.WriteString("\n=== Hasil lengkap disimpan di tests/output.txt ===\n")
}
