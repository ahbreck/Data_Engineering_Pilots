package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelvins/geocoder"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func saveJSON(filename string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func loadJSON(filename string, v any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func main() {

	startTime := time.Now() // start timer

	fmt.Println("Running func stream processing of crashes by zip code (no channels)")

	// Load geocoder API key
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	geocoder.ApiKey = os.Getenv("API_KEY")

	// Open CSV dataset
	file, err := os.Open("Traffic_Crashes_Mini_Dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Limiting rows for initial testing; change to len(data) for full dataset
	maxRows := len(data)

	cacheFile := "crash_in_zip_code_mp.json"
	crashInZip := make(map[string]int)

	geocodeCacheFile := "geocode_cache.json"
	geocodeCache := make(map[string]string)

	if fileExists(geocodeCacheFile) {
		err := loadJSON(geocodeCacheFile, &geocodeCache)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Loaded geocode cache.")
	} else {
		fmt.Println("No geocode cache found.")
	}

	fmt.Println("Processing traffic crashes...")

	newGeocodeEntries := 0
	for i, row := range data {
		if i == 0 || i > maxRows { // skip header and respect maxRows
			continue
		}

		latStr := row[46]
		longStr := row[47]

		if latStr == "" || longStr == "" {
			continue
		}

		key := fmt.Sprintf("%s,%s", latStr, longStr)
		var zip string

		if z, ok := geocodeCache[key]; ok {
			zip = z // reuse cached ZIP
		} else {
			lat, _ := strconv.ParseFloat(latStr, 64)
			long, _ := strconv.ParseFloat(longStr, 64)

			location := geocoder.Location{Latitude: lat, Longitude: long}
			addressList, err := geocoder.GeocodingReverse(location)
			if err != nil || len(addressList) == 0 {
				fmt.Printf("No results found for row %d\n", i)
				continue
			}

			zip = addressList[0].PostalCode
			if zip == "" {
				continue
			}

			// Add new entry to geocode cache
			geocodeCache[key] = zip
			newGeocodeEntries++
		}

		// Increment crash count for this ZIP
		crashInZip[zip]++

		// Print progress every 10 rows
		if i%10 == 0 {
			fmt.Printf("Processed %d rows... (new geocodes: %d)\n", i, newGeocodeEntries)
		}
	}

	// Save caches after processing
	err = saveJSON(cacheFile, crashInZip)
	if err != nil {
		log.Fatal(err)
	}

	err = saveJSON(geocodeCacheFile, geocodeCache)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nNumber of crashes per ZIP code:")
	for zip, count := range crashInZip {
		fmt.Printf("%s: %d\n", zip, count)
	}

	fmt.Printf("\nSaved %d new geocoding entries to %s\n", newGeocodeEntries, geocodeCacheFile)
	fmt.Printf("ZIP code crash counts saved to %s\n", cacheFile)

	elapsed := time.Since(startTime)
	fmt.Printf("\nTotal time elapsed: %s\n", elapsed)
}
