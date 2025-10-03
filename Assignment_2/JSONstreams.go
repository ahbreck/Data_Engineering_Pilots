package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Data struct {
	Key string `json:"key"`
	Val int    `json:"value"`
}

var DataRecords []Data

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

var MIN = 0
var MAX = 26

func getString(l int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return temp
}

// function to get the memory allocation currently in active use, in kilobytes
func getMemAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc / 1024 // return KB
}

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(e *json.Decoder, slice interface{}) error {
	return e.Decode(slice)
}

// Serialize serializes a slice with JSON records
func Serialize(e *json.Encoder, slice interface{}) error {
	return e.Encode(slice)
}

func main() {

	// sizes to test
	sizes := []int{10_000, 100_000, 1_000_000}

	// Use a for loop to repeat the experiment for each of the conditions in sizes slice
	for _, n := range sizes {
		fmt.Printf("\n=== %d records ===\n", n)

		// Garbage collection to ensure accuracy before measuring baseline memory
		runtime.GC()

		base_Mem := getMemAlloc()

		// Create sample data
		var i int
		var t Data
		for i = 0; i < 10000; i++ {
			t = Data{
				Key: getString(5),
				Val: random(1, 100),
			}
			DataRecords = append(DataRecords, t)
		}

		after_data_gen_Mem := getMemAlloc()

		// SERIALIZE
		// bytes.Buffer is both an io.Reader and io.Writer
		buf := new(bytes.Buffer)

		encoder := json.NewEncoder(buf)

		// Create check point to measure serialization time
		start := time.Now()

		err := Serialize(encoder, DataRecords)
		if err != nil {
			fmt.Println(err)
			return
		}

		// record time elapsed during serialization
		serDuration := time.Since(start)
		// fmt.Print("After Serialize:", buf)

		after_ser_Mem := getMemAlloc()

		// DESERIALIZE
		decoder := json.NewDecoder(buf)
		var temp []Data

		// Create check point to measure deserialization time
		start = time.Now()

		err = DeSerialize(decoder, &temp)

		// record time elapsed during serialization
		deserDuration := time.Since(start)

		after_deser_Mem := getMemAlloc()

		// print results
		fmt.Printf("Serialization time:   %v\n", serDuration)
		fmt.Printf("Deserialization time: %v\n", deserDuration)
		fmt.Printf("Baseline memory use:   %d KB\n", base_Mem)
		fmt.Printf("Memory use after data creation:   %d KB\n", after_data_gen_Mem)
		fmt.Printf("Memory use after serialization:   %d KB\n", after_ser_Mem)
		fmt.Printf("Memory use after deserialization:   %d KB\n", after_deser_Mem)

		fmt.Println()
	}
}
