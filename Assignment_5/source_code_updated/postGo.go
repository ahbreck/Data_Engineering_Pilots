package main

import (
	"fmt"

	//"github.com/mactsouk/post05"
	"github.com/ahbreck/Data_Engineering_Pilots/Assignment_5/source_code_updated/post05-main/post05-main"
	//"main.go/post05-main/post05-main"
)

func main() {
	post05.Hostname = "localhost"
	post05.Port = 5433
	post05.Username = "postgres"
	post05.Password = "root"
	post05.Database = "go"

	courses := []post05.MSDSCourse{
		{CID: "MSDS500", CNAME: "Data Science Foundations", CPREREQ: "None"},
		{CID: "MSDS510", CNAME: "Applied Statistics", CPREREQ: "MSDS500"},
		{CID: "MSDS520", CNAME: "Machine Learning", CPREREQ: "MSDS510"},
		{CID: "MSDS530", CNAME: "Data Engineering", CPREREQ: "MSDS500"},
		{CID: "MSDS540", CNAME: "Data Visualization", CPREREQ: "MSDS500"},
	}

	if err := post05.SeedCourseCatalog(courses); err != nil {
		fmt.Println("unable to seed course catalog:", err)
		return
	}

}
