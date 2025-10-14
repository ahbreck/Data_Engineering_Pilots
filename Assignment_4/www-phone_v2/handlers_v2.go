package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT = ":1234"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", Body)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get course ID
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)

	courseID := paramStr[0]
	err := deleteEntry(courseID)
	if err != nil {
		fmt.Println(err)
		Body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", Body)
		return
	}

	Body := courseID + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", Body)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := list()
	fmt.Fprintf(w, "%s", Body)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := fmt.Sprintf("Total entries: %d\n", len(data))
	fmt.Fprintf(w, "%s", Body)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	// Split URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}

	courseID := paramStr[0]
	course_name := paramStr[1]
	course_prereqs := paramStr[2]

	t := strings.ReplaceAll(courseID, "-", "")
	if !matchTel(t) {
		fmt.Println("Not a valid course ID:", courseID)
		return
	}

	temp := &MSDSCourse{CID: courseID, CNAME: course_name, CPREREQ: course_prereqs}
	err := insert(temp)

	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		Body := "Failed to add record\n"
		fmt.Fprintf(w, "%s", Body)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		Body := "New record added successfully\n"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", Body)
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get Search value from URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	var Body string
	courseID := paramStr[0]
	t := search(courseID)
	if t == nil {
		w.WriteHeader(http.StatusNotFound)
		Body = "Could not be found: " + courseID + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		Body = t.CID + " " + t.CNAME + " " + t.CPREREQ + "\n"
	}

	fmt.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", Body)
}
