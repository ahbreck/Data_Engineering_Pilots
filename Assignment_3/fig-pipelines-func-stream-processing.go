package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Create a logical condition to use in the functions with Go generics to allow any type with an underlying type of int or float64
type Number interface {
	~int | ~float64
}

// The "Number" logical condition was defined in the 'fig-functional-pipeline-combination.go' script in this same package.
// it allows either integer or float64
func multiply_v2[T Number](value, multiplier T) T {
	return value * multiplier
}

func add_v2[T Number](value, additive T) T {
	return value + additive
}

// create functions to generate integer slices and float slices

func make_int_slice(length int) []int {
	ints := make([]int, length)
	for i := range ints {
		ints[i] = rand.Intn(10)
	}
	return ints
}

func make_float_slice(length int) []float64 {
	floats := make([]float64, length)
	for i := range floats {
		floats[i] = rand.Float64() * 10
	}
	return floats
}

// Result struct to store timing data for later summary
type Result struct {
	SampleSize  int
	Trial       int
	IntegerTime time.Duration
	FloatTime   time.Duration
}

func main() {

	samples := []int{10_000, 100_000, 1_000_000}

	// this variable will store the results. it is a "Result" struct, defined above
	var results []Result

	for _, size := range samples {

		fmt.Printf("\n=== Experiments with samples of size: %d ===\n", size)

		for i := 1; i < 3; i++ {
			fmt.Printf("\nTrial number %d \n", i)

			ints := make_int_slice(size)

			// Create check point to measure time for integers
			start := time.Now()
			for _, v := range ints {
				fmt.Println(multiply_v2(add_v2(multiply_v2(v, 2), 1), 2))
			}

			// record time elapsed for integers
			intDuration := time.Since(start)

			floats := make_float_slice(size)

			// Create check point to measure time for floats
			start = time.Now()

			for _, v := range floats {
				fmt.Println(multiply_v2(add_v2(multiply_v2(v, 2), 1), 2))
			}

			// record time elapsed for integers
			floatDuration := time.Since(start)

			// Store results
			results = append(results, Result{
				SampleSize:  size,
				Trial:       i,
				IntegerTime: intDuration,
				FloatTime:   floatDuration,
			})

		}

	}

	// Print summary table of results at the end
	fmt.Printf("\n\n======= SUMMARY RESULTS FOR STREAM PROCESSING =======\n")
	fmt.Printf("%-12s %-8s %-15s %-15s\n", "SampleSize", "Trial", "IntegerTime", "FloatTime")
	fmt.Println("---------------------------------------------------------")
	for _, r := range results {
		fmt.Printf("%-12d %-8d %-15v %-15v\n",
			r.SampleSize, r.Trial, r.IntegerTime, r.FloatTime)
	}

}
