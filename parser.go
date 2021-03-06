package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

// Vacancy is a type with the fields used for parsing
type Vacancy struct {
	Position string
	Company  string
	Link     string
	City     string
	Details  string
}

// PrepareDB and return it
func PrepareDB() *sql.DB {
	database, err := sql.Open("sqlite3", "vacancies.sqlite")
	if err != nil {
		panic(err)
	}

	return database
}

// CreateTable if not exists
func CreateTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS vacancies(
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Position TEXT,
		Company TEXT,
		Link TEXT,
		City TEXT,
		Details TEXT
	);
	`

	_, err := db.Exec(sqlTable)
	if err != nil {
		panic(err)
	}
}

// InsertVacancy stores found vacancis in the db
func InsertVacancy(db *sql.DB, job Vacancy) {
	sqlAdd := `
	INSERT INTO vacancies(Position, Company, Link, City, Details) VALUES (?, ?, ?, ?, ?);
	`

	stmt, err := db.Prepare(sqlAdd)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	fmt.Println("Inserting " + job.Position)
	_, err2 := stmt.Exec(job.Position, job.Company, job.Link, job.City, job.Details)
	if err2 != nil {
		panic(err2)
	}

}

// ExistsVacancy checks if given job exists in the database
func ExistsVacancy(db *sql.DB, job Vacancy) bool {
	sqlRead := fmt.Sprintf(`SELECT ID FROM vacancies WHERE Position="%s" AND Company="%s" AND City="%s";`,
		job.Position, job.Company, job.City)

	// fmt.Println(sqlRead)
	rows, err := db.Query(sqlRead)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
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
func ExampleScrape() []Vacancy {
	// Request the HTML page.
	res, err := http.Get("https://hh.ru/search/vacancy?text=%28junior+OR+trainee+OR+intern%29+and+%28" +
		"Go+OR+Golang+OR+Python%29&only_with_salary=false&order_by=publication_time&specialization=1" +
		"&area=1&enable_snippets=true&clusters=true&experience=noExperience&salary=")
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

	jobs := []Vacancy{}

	// buttons := doc.Find("a.bloko-button")

	buttons := doc.Find("a.bloko-button").FilterFunction(func(i int, s *goquery.Selection) bool {
		_, exists := s.Attr("data-page")
		fmt.Println(exists)
		return exists
	})
	// numberButtons := buttons.Size()
	lastPage := buttons.Last().Text()

	_, err2 := strconv.ParseInt(lastPage, 10, 0)
	if err2 != nil {
		fmt.Println("Button wasnt numeric ", lastPage)

		lastPage = buttons.Eq(-2).Text()
		fmt.Println("Last button now ", lastPage)
		_, err3 := strconv.ParseInt(lastPage, 10, 0)
		if err3 != nil {
			fmt.Println(lastPage)
			panic("Unknown number of pages")
		}
	}

	buttons.Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	// fmt.Println(lastPage)

	// Find the review items
	doc.Find(".vacancy-serp-item").Each(func(i int, s *goquery.Selection) {
		address := s.Find("span.vacancy-serp-item__meta-info").First()
		city := strings.Split(address.Text(), ",")[0]
		// fmt.Println(city)

		position, link, employer := "", "", ""

		links := s.Find("a")
		for j := range links.Nodes {
			single := links.Eq(j)
			data, _ := single.Attr("data-qa")

			if data == "vacancy-serp__vacancy-title" {
				position = single.Text()
				link, _ = single.Attr("href")
				// fmt.Println(position, link)
			}

			if data == "vacancy-serp__vacancy-employer" {
				employer = single.Text()
				// fmt.Println(employer)
			}
		}

		details := "DESCRIPTION GOES HERE"

		job := Vacancy{position, employer, link, city, details}
		jobs = append(jobs, job)
		// fmt.Println()

	})

	return jobs
}

func main() {
	db := PrepareDB()
	CreateTable(db)

	allJobs := ExampleScrape()
	// link := `https://hh.ru/search/vacancy?text=%28junior+OR+trainee+OR+intern%29+and+%28" +
	// 		"Go+OR+Golang+OR+Python%29&only_with_salary=false&order_by=publication_time&specialization=1" +
	// 		"&area=113&enable_snippets=true&clusters=true&experience=noExperience&salary=`

	for _, job := range allJobs {
		if ExistsVacancy(db, job) {
			continue
		} else {
			fmt.Println("Found new job posting!")
			new := fmt.Sprintf("%s\n %s\n %s\n %s\n", job.Position, job.Company, job.City, job.Link)
			fmt.Println(new)

			InsertVacancy(db, job)
		}

	}
	// fmt.Println(isExists)

	// fmt.Println(allJobs)
	fmt.Println("Finished")
}
