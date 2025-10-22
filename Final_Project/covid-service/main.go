package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type CovidRecord struct {
	ZIP                   string  `json:"zip_code"`
	WeekStart             string  `json:"week_start"`
	WeekEnd               string  `json:"week_end"`
	CaseRateWeekly        float64 `json:"case_rate_weekly,string"`
	PercentTestedPositive float64 `json:"percent_tested_positive_weekly,string"`
}

const dayOffset = 2029 // fixed offset to simulate "current" data

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	apiURL := os.Getenv("COVID_API_URL")
	fetchIntervalDaysStr := os.Getenv("FETCH_INTERVAL_DAYS")
	runContinuously := os.Getenv("RUN_CONTINUOUSLY")

	fetchIntervalDays, err := strconv.Atoi(fetchIntervalDaysStr)
	if err != nil || fetchIntervalDays <= 0 {
		fetchIntervalDays = 7
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %v", err)
	}
	defer db.Close()

	createTableIfMissing(db)

	for {
		err := fetchAndInsertData(apiURL, db)
		if err != nil {
			log.Printf("Error in fetch/insert: %v", err)
		}

		if runContinuously != "true" {
			log.Println("RUN_CONTINUOUSLY=false → exiting after one run.")
			break
		}

		log.Printf("Sleeping for %d days before next fetch...\n", fetchIntervalDays)
		time.Sleep(time.Duration(fetchIntervalDays) * 24 * time.Hour)
	}
}

func createTableIfMissing(db *sql.DB) {
	ddl := `
	CREATE TABLE IF NOT EXISTS covid_zip_weekly (
		zip_code CHAR(9),
		week_start DATE,
		week_end DATE,
		case_rate_weekly NUMERIC,
		percent_tested_positive_weekly NUMERIC,
		records_fetched_at TIMESTAMP DEFAULT now(),
		PRIMARY KEY (zip_code, week_start)
	);`
	_, err := db.Exec(ddl)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func fetchAndInsertData(apiURL string, db *sql.DB) error {
	fullURL := fmt.Sprintf("%s?$limit=50000", apiURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("fetch error: %v", err)
	}
	defer resp.Body.Close()

	var records []CovidRecord
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return fmt.Errorf("decode error: %v", err)
	}

	for _, r := range records {
		if r.ZIP == "" || r.WeekStart == "" || r.WeekEnd == "" {
			continue
		}

		startDate, err1 := time.Parse("2006-01-02T15:04:05.000", r.WeekStart)
		if err1 != nil {
			startDate, _ = time.Parse("2006-01-02", r.WeekStart)
		}
		endDate, err2 := time.Parse("2006-01-02T15:04:05.000", r.WeekEnd)
		if err2 != nil {
			endDate, _ = time.Parse("2006-01-02", r.WeekEnd)
		}

		// Shift by fixed offset
		startShifted := startDate.AddDate(0, 0, dayOffset)
		endShifted := endDate.AddDate(0, 0, dayOffset)

		_, err = db.Exec(`
			INSERT INTO covid_zip_weekly 
				(zip_code, week_start, week_end, case_rate_weekly, percent_tested_positive_weekly) 
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (zip_code, week_start) DO UPDATE
				SET case_rate_weekly = EXCLUDED.case_rate_weekly,
				    percent_tested_positive_weekly = EXCLUDED.percent_tested_positive_weekly;
		`, r.ZIP, startShifted, endShifted, r.CaseRateWeekly, r.PercentTestedPositive)
		if err != nil {
			log.Printf("Insert failed for ZIP %s week %s: %v", r.ZIP, r.WeekStart, err)
		}
	}
	log.Printf("Inserted or updated %d records.", len(records))
	return nil
}
