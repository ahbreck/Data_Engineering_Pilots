package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelvins/geocoder"
)

// ----------------------------
// Helper functions
// ----------------------------
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

// ----------------------------
// Worker functions
// ----------------------------

// rowGenerator sends rows downstream for processing
func rowGenerator(done <-chan struct{}, data [][]string, maxRows int) <-chan []string {
	out := make(chan []string)
	go func() {
		defer close(out)
		for i, row := range data {
			if i == 0 || i > maxRows {
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

// geocodeWorker converts lat/long rows to ZIP codes using geocode cache
func geocodeWorker(done <-chan struct{}, rows <-chan []string, geocodeCache map[string]string, mu *sync.Mutex) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for row := range rows {
			latStr, longStr := row[46], row[47]
			if latStr == "" || longStr == "" {
				continue
			}

			key := fmt.Sprintf("%s,%s", latStr, longStr)

			var zip string
			mu.Lock()
			z, found := geocodeCache[key]
			mu.Unlock()

			if found {
				zip = z
			} else {
				lat, _ := strconv.ParseFloat(latStr, 64)
				long, _ := strconv.ParseFloat(longStr, 64)
				location := geocoder.Location{Latitude: lat, Longitude: long}
				addressList, err := geocoder.GeocodingReverse(location)
				if err != nil || len(addressList) == 0 {
					continue
				}
				zip = addressList[0].PostalCode
				if zip == "" {
					continue
				}
				mu.Lock()
				geocodeCache[key] = zip
				mu.Unlock()
			}

			select {
			case <-done:
				return
			case out <- zip:
			}
		}
	}()
	return out
}

// mergeChannels merges multiple zip channels into one
func mergeChannels(done <-chan struct{}, chans ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		defer wg.Done()
		for val := range c {
			select {
			case <-done:
				return
			case out <- val:
			}
		}
	}

	wg.Add(len(chans))
	for _, c := range chans {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// countZips tallies ZIP code occurrences
func countZips(done <-chan struct{}, zips <-chan string) map[string]int {
	counts := make(map[string]int)
	i := 0
	for zip := range zips {
		counts[zip]++
		i++
		if i%20 == 0 {
			fmt.Printf("Processed %d ZIPs so far...\n", i)
		}
	}
	return counts
}

// ----------------------------
// Main
// ----------------------------
func main() {
	start := time.Now()

	// Load .env and set API key
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	geocoder.ApiKey = os.Getenv("API_KEY")

	// Load dataset
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

	maxRows := len(data)

	// Setup cancellation channel
	done := make(chan struct{})
	defer close(done)

	// Load geocode cache
	geocodeCacheFile := "geocode_cache.json"
	geocodeCache := make(map[string]string)
	if fileExists(geocodeCacheFile) {
		_ = loadJSON(geocodeCacheFile, &geocodeCache)
		fmt.Printf("Loaded %d geocoded entries from cache.\n", len(geocodeCache))
	}

	// Start pipeline
	rowChan := rowGenerator(done, data, maxRows)

	// Launch multiple geocode workers
	numWorkers := 5
	var mu sync.Mutex
	var workerChans []<-chan string
	for i := 0; i < numWorkers; i++ {
		workerChans = append(workerChans, geocodeWorker(done, rowChan, geocodeCache, &mu))
	}

	// Merge worker outputs
	zipChan := mergeChannels(done, workerChans...)

	// Count ZIP codes
	counts := countZips(done, zipChan)

	// Save caches
	_ = saveJSON("crash_in_zip_code_mp.json", counts)
	_ = saveJSON(geocodeCacheFile, geocodeCache)

	// Print summary
	fmt.Println("\nNumber of crashes per ZIP code:")
	for zip, c := range counts {
		fmt.Printf("%s: %d\n", zip, c)
	}

	fmt.Printf("\nTotal execution time: %s\n", time.Since(start))
}
