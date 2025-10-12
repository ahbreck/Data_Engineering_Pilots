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

// Helper functions

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

// Channel to stream rows from the CSV into the pipeline
func rowGenerator(done <-chan struct{}, data [][]string, maxRows int) <-chan []string {
	out := make(chan []string)
	go func() {
		defer close(out)
		for i, row := range data {
			if i == 0 || i > maxRows { // skip header and respect maxRows, if applicable
				continue
			}
			select {
			case <-done:
				return
			case out <- row:
			}
		}
	}()
	return out
}

// Channel to geocode zip codes (or load from cache if previously mapped in a prior run) and send results on
func geocodeWorker(done <-chan struct{}, in <-chan []string, geocodeCache map[string]string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		count := 0
		for row := range in {
			select {
			case <-done:
				return
			default:
			}

			latStr := row[46]
			longStr := row[47]
			if latStr == "" || longStr == "" {
				continue
			}
			key := fmt.Sprintf("%s,%s", latStr, longStr)
			var zip string
			if z, ok := geocodeCache[key]; ok {
				zip = z
			} else {
				lat, _ := strconv.ParseFloat(latStr, 64)
				long, _ := strconv.ParseFloat(longStr, 64)
				location := geocoder.Location{Latitude: lat, Longitude: long}
				addressList, err := geocoder.GeocodingReverse(location)
				if err != nil || len(addressList) == 0 {
					fmt.Printf("No results for %s\n", key)
					continue
				}
				zip = addressList[0].PostalCode
				if zip == "" {
					continue
				}
				geocodeCache[key] = zip
			}

			select {
			case <-done:
				return
			case out <- zip:
			}

			count++
			if count%10 == 0 { // progress every 10 processed rows
				fmt.Printf("Geocoded %d rows...\n", count)
			}
		}
	}()
	return out
}

// Channel to consume the zip code data from geocodeWorker channel and produce counts of crashes per zip code
func countZips(done <-chan struct{}, in <-chan string) map[string]int {
	counts := make(map[string]int)
	for zip := range in {
		select {
		case <-done:
			return counts
		default:
		}
		counts[zip]++
	}
	return counts
}

// ----------------------------
// Main
// ----------------------------
func main() {

	startTime := time.Now() // start timer

	fmt.Println("Running channel stream processing of crashes by zip code")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	geocoder.ApiKey = os.Getenv("API_KEY")

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

	maxRows := len(data) // adjust to len(data) after testing.
	done := make(chan struct{})
	defer close(done)

	// ----------------------
	// Load geocode cache
	// ----------------------
	geocodeCacheFile := "geocode_cache.json"
	geocodeCache := make(map[string]string)
	if fileExists(geocodeCacheFile) {
		if err := loadJSON(geocodeCacheFile, &geocodeCache); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Loaded geocode cache.")
	} else {
		fmt.Println("No geocode cache found.")
	}

	// ----------------------
	// Pipeline
	// ----------------------
	rows := rowGenerator(done, data, maxRows)
	zips := geocodeWorker(done, rows, geocodeCache)
	crashCounts := countZips(done, zips)

	// ----------------------
	// Save caches
	// ----------------------
	err = saveJSON("crash_in_zip_code_mp.json", crashCounts)
	if err != nil {
		log.Fatal(err)
	}
	err = saveJSON(geocodeCacheFile, geocodeCache)
	if err != nil {
		log.Fatal(err)
	}

	// ----------------------
	// Print summary
	// ----------------------
	fmt.Println("\nNumber of crashes per ZIP code:")
	for zip, count := range crashCounts {
		fmt.Printf("%s: %d\n", zip, count)
	}
	fmt.Println("\nPipeline complete.")
	elapsed := time.Since(startTime)
	fmt.Printf("\nTotal time elapsed: %s\n", elapsed)
}
