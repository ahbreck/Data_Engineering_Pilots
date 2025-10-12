package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Create a logical condition to use in the functions with Go generics to allow any type with an underlying type of int or float64
type Number interface {
	~int | ~float64
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

func generator[T Number](done <-chan interface{}, numbers ...T) <-chan T {
	numStream := make(chan T)
	go func() {
		defer close(numStream)
		for _, i := range numbers {
			select {
			case <-done:
				return
			case numStream <- i:
			}
		}
	}()
	return numStream
}

func multiplyWorker[T Number](done <-chan interface{}, in <-chan T, factor T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v * factor:
			}
		}
	}()
	return out
}

func addWorker[T Number](done <-chan interface{}, in <-chan T, offset T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v + offset:
			}
		}
	}()
	return out
}

// --- Fan-in helper: merges multiple worker output channels
func fanIn[T Number](done <-chan interface{}, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()
		for v := range c {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func multiWorkerPipeline[T Number](done <-chan interface{}, input <-chan T, workers int) <-chan T {
	// Step 1: Fan-out multiply
	multiplied := make([]<-chan T, workers)
	for i := 0; i < workers; i++ {
		multiplied[i] = multiplyWorker(done, input, 2)
	}
	merged1 := fanIn(done, multiplied...)

	// Step 2: Fan-out add
	added := make([]<-chan T, workers)
	for i := 0; i < workers; i++ {
		added[i] = addWorker(done, merged1, 1)
	}
	merged2 := fanIn(done, added...)

	// Step 3: Fan-out multiply again
	final := make([]<-chan T, workers)
	for i := 0; i < workers; i++ {
		final[i] = multiplyWorker(done, merged2, 2)
	}

	return fanIn(done, final...)
}

func main() {

	samples := []int{10_000, 100_000, 1_000_000}

	// this variable will store the results. it is a "Result" struct, defined above
	var results []Result

	workers := 4

	for _, size := range samples {

		fmt.Printf("\n=== Experiments with samples of size: %d ===\n", size)

		for i := 1; i < 3; i++ {
			fmt.Printf("\nTrial number %d \n", i)

			ints := make_int_slice(size)

			// Create check point to measure time for integers
			start := time.Now()

			done := make(chan interface{})

			intStream := generator(done, ints...)
			intpipeline := multiWorkerPipeline(done, intStream, workers)
			for v := range intpipeline {
				fmt.Println(v)
			}

			// record time elapsed for integers
			intDuration := time.Since(start)
			close(done)

			floats := make_float_slice(size)

			// Create check point to measure time for floats
			start = time.Now()

			done = make(chan interface{})

			floatStream := generator(done, floats...)
			floatpipeline := multiWorkerPipeline(done, floatStream, workers)
			fmt.Println(v)
		}

		// record time elapsed for integers
		floatDuration := time.Since(start)

		close(done)

		// Store results
		results = append(results, Result{
			SampleSize:  size,
			Trial:       i,
			IntegerTime: intDuration,
			FloatTime:   floatDuration,
		})

	}

	// Print summary table of results at the end
	fmt.Printf("\n\n======= SUMMARY RESULTS FOR MULTIWORKER CHANNEL STREAM PROCESSING =======\n")
	fmt.Printf("%-12s %-8s %-15s %-15s\n", "SampleSize", "Trial", "IntegerTime", "FloatTime")
	fmt.Println("---------------------------------------------------------")
	for _, r := range results {
		fmt.Printf("%-12d %-8d %-15v %-15v\n",
			r.SampleSize, r.Trial, r.IntegerTime, r.FloatTime)
	}

}
