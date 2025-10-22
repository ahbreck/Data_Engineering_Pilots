package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Trip struct {
	TripID    string `json:"trip_id"`
	StartTS   string `json:"trip_start_timestamp"`
	EndTS     string `json:"trip_end_timestamp"`
	TripMiles string `json:"trip_miles"`
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %v", err)
	}
	defer db.Close()

	// Example: Pull first 5 records from Taxi Trips dataset
	resp, err := http.Get("https://data.cityofchicago.org/resource/wrvz-psew.json?$limit=5")
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	var trips []Trip
	if err := json.NewDecoder(resp.Body).Decode(&trips); err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	for _, t := range trips {
		_, err := db.Exec(`
			INSERT INTO trips (trip_id, trip_start, trip_end, trip_miles)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (trip_id) DO NOTHING;
		`, t.TripID, t.StartTS, t.EndTS, t.TripMiles)
		if err != nil {
			log.Printf("Insert failed for trip %s: %v", t.TripID, err)
		}
	}

	fmt.Println("Inserted sample trips successfully.")
}
