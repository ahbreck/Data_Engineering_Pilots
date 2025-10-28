package main

import (
	"fmt"
)

func main() {
	Hostname = "localhost"
	Port = 5435
	Username = "postgres"
	Password = "root"
	Database = "go"

	courses := []MSDSCourse{
		{CID: "MSDS420", CNAME: "Database Systems", CPREREQ: "None"},
		{CID: "MSDS458", CNAME: "Artificial Intelligence and Deep Learning", CPREREQ: "MSDS420 and MSDS422"},
		{CID: "MSDS460", CNAME: "Decision Analytics", CPREREQ: "MSDS400 and MSDS401"},
		{CID: "MSDS422", CNAME: "Practical Machine Learning", CPREREQ: "MSDS400 and MSDS401"},
		{CID: "MSDS432", CNAME: "Foundations Of Data Engineering", CPREREQ: "MSDS420"},
	}

	if err := SeedCourseCatalog(courses); err != nil {
		fmt.Println("unable to seed course catalog:", err)
		return
	}

}
