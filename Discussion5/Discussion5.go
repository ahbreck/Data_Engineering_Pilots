package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelvins/geocoder"
)

type TripData struct {
	TripID                  string           `json:"trip_id"`
	PickupCentroidLatitude  string           `json:"pickup_centroid_latitude"`
	PickupCentroidLongitude string           `json:"pickup_centroid_longitude"`
	PickupLocation          string           `json:"pickup_location"`
	Number                  int              `json:"number"`
	Street                  string           `json:"street"`
	City                    string           `json:"city"`
	State                   string           `json:"state"`
	Zipcode                 string           `json:"zipcode"`
	Address                 geocoder.Address `json:"full_returned_address"`
}

func main() {

	// open  csv dataset
	dataset, err := os.Open("Taxi_Trips_Sample.csv")

	if err != nil {
		log.Fatal(err)
	}

	defer dataset.Close()

	// read data using csv.Reader
	csvReader := csv.NewReader(dataset)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// converts data to map with trip_id as key and taxi_data struct as values
	trips_mp := createTripsMap(data)

	fmt.Println("---------------------Sample Results for Reverse Geocoding of Taxi Pickups-----------------")

	sample, err := json.MarshalIndent(trips_mp, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(sample))
	fmt.Printf("\n\n")

}

func createTripsMap(data [][]string) map[string]TripData {

	trips_mp := make(map[string]TripData)

	fmt.Println("Creating Map from Trip Data")

	// Load .env file to get geocoder API key
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	geocoder.ApiKey = os.Getenv("API_KEY")

	// Just doing two records for now
	for i := 1; i < 3; i++ {

		//initializing the list
		var TripRecord TripData

		TripRecord.TripID = data[i][0]

		if data[i][17] == "" || data[i][18] == "" {
			continue
		} else {
			TripRecord.PickupCentroidLatitude = data[i][17]
			TripRecord.PickupCentroidLongitude = data[i][18]
			TripRecord.PickupLocation = data[i][19]
		}

		pickup_latitude_float, _ := strconv.ParseFloat(TripRecord.PickupCentroidLatitude, 64)
		pickup_longitude_float, _ := strconv.ParseFloat(TripRecord.PickupCentroidLongitude, 64)

		location := geocoder.Location{
			Latitude:  pickup_latitude_float,
			Longitude: pickup_longitude_float,
		}

		address_list, _ := geocoder.GeocodingReverse(location)

		// Ignoring the entry if location of the data is not available
		if len(address_list) == 0 {
			fmt.Printf("No results found for trip pickup at latitude : %f and longitude : %f \n", pickup_latitude_float, pickup_longitude_float)
			continue
		}

		address := address_list[0]

		TripRecord.Number = address.Number
		TripRecord.Street = address.Street
		TripRecord.City = address.City
		TripRecord.State = address.State
		TripRecord.Zipcode = address.PostalCode
		TripRecord.Address = address

		trips_mp[data[i][0]] = TripRecord

	}

	return trips_mp

}
