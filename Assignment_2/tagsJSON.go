package main

import (
	"encoding/json"
	"fmt"
)

type MSDSCourse struct {
	CID     string `json:"courseI_D"`
	CNAME   string `json:"course_name"`
	CPREREQ string `json:"prerequisite"`
}

func main() {

	// Create some courses
	c1 := MSDSCourse{CID: "MSDS432", CNAME: "Foundations of Data Engineering", CPREREQ: "MSDS400 and MSDS420"}
	c2 := MSDSCourse{CID: "MSDS420", CNAME: "Database Systems", CPREREQ: "None"}
	c3 := MSDSCourse{CID: "MSDS458", CNAME: "Artificial Intelligence and Deep Learning", CPREREQ: "MSDS420 and MSDS422"}
	c4 := MSDSCourse{CID: "MSDS460", CNAME: "Decision Analytics", CPREREQ: "MSDS400 and MSDS401"}
	c5 := MSDSCourse{CID: "MSDS422", CNAME: "Practical Machine Learning", CPREREQ: "MSDS400 and MSDS401"}

	// Create slice object to store courses
	courses_slice := []MSDSCourse{c1, c2, c3, c4, c5}

	// Create array object to store courses
	courses_array := [5]MSDSCourse{c1, c2, c3, c4, c5}

	// Create map object to store courses
	courses_map := make(map[string]MSDSCourse)
	courses_map[c1.CID] = c1
	courses_map[c2.CID] = c2
	courses_map[c3.CID] = c3
	courses_map[c4.CID] = c4
	courses_map[c5.CID] = c5

	// Slice example
	slice_bytes, err := json.Marshal(courses_slice)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println()
		fmt.Printf("Slice decoded with value %s\n", slice_bytes)
		fmt.Println()
		fmt.Printf("Original format was %s\n", courses_slice)
		fmt.Println()
	}

	// Array example
	array_bytes, err := json.Marshal(courses_array)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Array decoded with value %s\n", array_bytes)
		fmt.Println()
		fmt.Printf("Original format was %s\n", courses_array)
		fmt.Println()
	}

	// Map example
	map_bytes, err := json.Marshal(courses_map)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Map decoded with value %s\n", map_bytes)
		fmt.Println()
		fmt.Printf("Original format was %s\n", courses_map)
		fmt.Println()
	}
}
