package main

import (
	"fmt"
	"math/rand"
)

func main() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}
	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	//ints := []int{1, 2, 3, 4}

	ints := make([]int, 100)
	for i := range ints {
		ints[i] = rand.Intn(10)
	}

	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}
