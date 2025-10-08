package main

import (
	"fmt"
	"math/rand"
)

// Create a logical condition to use in the functions with Go generics to allow any type with an underlying type of int or float64
type Number interface {
	~int | ~float64
}

// Moved the function declarations to the package area to allow use of Go generics so the functions can be more flexible to use either int or float64 inputs
// Using Go version 1.24.4, which allows for this
func multiply[T Number](values []T, multiplier T) []T {
	multipliedValues := make([]T, len(values))
	for i, v := range values {
		multipliedValues[i] = v * multiplier
	}
	return multipliedValues
}

func add[T Number](values []T, additive T) []T {
	addedValues := make([]T, len(values))
	for i, v := range values {
		addedValues[i] = v + additive
	}
	return addedValues
}

// create functions to generate integer slices and float slices

func main() {

	ints := make([]int, 100)
	for i := range ints {
		ints[i] = rand.Intn(10)
	}

	floats := make([]float64, 100)
	for i := range floats {
		floats[i] = rand.Float64()
	}

	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

	for _, v := range add(multiply(floats, 2), 1) {
		fmt.Println(v)
	}

}
