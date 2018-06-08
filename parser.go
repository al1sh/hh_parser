package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Vacancy is a type with the fields used for parsing
type Vacancy struct {
	ID         int
	Position   string
	Company    string
	Link       string
	City       string
	Experience string
	Language   string
	Date       string
	Details    string
}

// PrepareDB and return it
func PrepareDB() *sql.DB {
	database, err := sql.Open("sqlite3", "./vacancies.db")
	if err != nil {
		panic(err)
	}

	return database
}

// CreateTable if not exists
func CreateTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS vacancies(
		Id TEXT NOT NULL ,
		Position TEXT,
		Company TEXT,
		Link TEXT,
		City TEXT,
		Experience TEXT,
		Language: TEXT,
		Date: TEXT
		DETAILS: TEXT
	);
	`

	_, err := db.Exec(sqlTable)
	if err != nil {
		panic(err)
	}
}

// InsertVacancy stores found vacancis in the db
func InsertVacancy(db *sql.DB, jobs []Vacancy) {
	sqlAdd := `
	INSERT OR REPLACE INTO vacancies(
		Position 
		Company 
		Link 
		City 
		Experience 
		Language 
		Date
		DETAILS
	) values(?, ?, ?, ?, ?, ? ,? ,?)
	`

	stmt, err := db.Prepare(sqlAdd)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for _, job := range jobs {
		_, err2 := stmt.Exec(job.Position, job.Company, job.Link, job.City, job.Experience, job.Language, job.Date, job.Details)
		if err2 != nil {
			panic(err2)
		}
	}
}

// ExistsVacancy checks if given job exists in the database
func ExistsVacancy(db *sql.DB) bool {
	sqlRead := `
	SELECT Position, Company, City, Experience, Language FROM vacancies
	`

	rows, err := db.Query(sqlRead)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows != nil {
		return true
	}
	return false

	// var result []Vacancy
	// for rows.Next() {
	// 	item := Vacancy{}
	// 	err2 := rows.Scan(&item.Id, &item.Name, &item.Phone)
	// 	if err2 != nil {
	// 		panic(err2)
	// 	}
	// 	result = append(result, item)
	// }
}

// ExampleScrape scrapes given URL
func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get("https://hh.ru/search/vacancy?text=%28junior+OR+trainee+OR+intern%29+and+%28" +
		"Go+OR+Golang+OR+Python%29&only_with_salary=false&order_by=publication_time&specialization=1" +
		"&area=113&enable_snippets=true&clusters=true&experience=noExperience&salary=")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".vacancy-serp-item").Each(func(i int, s *goquery.Selection) {
		address := s.Find("span.vacancy-serp-item__meta-info").First()
		city := strings.Split(address.Text(), ",")[0]
		fmt.Println(city)

		links := s.Find("a")
		for j := range links.Nodes {
			single := links.Eq(j)
			data, _ := single.Attr("data-qa")

			if data == "vacancy-serp__vacancy-title" {
				position := single.Text()
				link, _ := single.Attr("href")
				fmt.Println(position, link)
			}

			if data == "vacancy-serp__vacancy-employer" {
				employer := single.Text()
				fmt.Println(employer)
			}
		}
		fmt.Println()

	})
}

func main() {

	ExampleScrape()
	fmt.Println("Finished")
}
