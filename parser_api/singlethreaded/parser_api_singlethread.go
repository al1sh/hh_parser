package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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

func (v Vacancy) String() string {
	return fmt.Sprintf("Position: %s\n Company: %s\n City: %s\n Description: %s\n Link: %s\n", v.Position, v.Company, v.City, v.Details, v.Link)
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
	sqlRead, errPrep := db.Prepare(`SELECT ID FROM vacancies WHERE Position=? AND Company=? AND City=?`)
	if errPrep != nil {
		panic(errPrep)
	}
	// fmt.Println(sqlRead)
	rows, err := sqlRead.Query(job.Position, job.Company, job.City)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}
	return false
}

// ExampleScrape scrapes given URL
func ExampleScrape(url string) []Vacancy {
	// Request the HTML page.

	jobs := []Vacancy{}

	apiGet := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Personal Job Search. Feedback: firstrestrest@gmail.com")

	res, getErr := apiGet.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var data map[string]interface{}
	errParse := json.Unmarshal([]byte(body), &data)
	if errParse != nil {
		panic("Couldn't read in json")
	}

	items := data["items"].([]interface{})

	for _, item := range items {
		job := Vacancy{}
		info := item.(map[string]interface{})

		name := info["name"].(string)
		job.Position = name
		// fmt.Println(name)

		area := info["area"].(map[string]interface{})
		city := area["name"].(string)
		job.City = city
		// fmt.Println(city)

		link := info["alternate_url"].(string)
		job.Link = link
		// fmt.Println(link)

		employer := info["employer"].(map[string]interface{})
		company := employer["name"].(string)
		job.Company = company
		// fmt.Println(company)
		// fmt.Println(job)

		jobs = append(jobs, job)
	}

	return jobs

}

func main() {
	start := time.Now()
	db := PrepareDB()
	CreateTable(db)

	jobs := []string{"python", "Go+OR+Golang", "Project+manager+AND+English"}
	experiences := []string{"noExperience", "between1And3"}
	cities := []string{"1", "2", "1624"}
	allJobs := []Vacancy{}

	for _, job := range jobs {
		for _, experience := range experiences {
			for _, city := range cities {
				url := "https://api.hh.ru/vacancies?text=" + job + "&area=" + city + "&experience=" + experience + "&per_page=100&specialization=1"
				allJobs = append(ExampleScrape(url), allJobs...)
			}
		}
	}

	// allJobs := ExampleScrape()

	for _, position := range allJobs {
		// fmt.Println(position)
		if ExistsVacancy(db, position) {
			continue
		} else {
			fmt.Println("Found new job posting!")
			fmt.Println(position)

			InsertVacancy(db, position)
		}

	}

	elapsed := time.Since(start)
	fmt.Printf("Search took %s", elapsed)

}
