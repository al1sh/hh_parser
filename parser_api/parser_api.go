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
func ExampleScrape() {
	// Request the HTML page.
	url := "https://api.hh.ru/vacancies?text=java&area=2"

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

	pages := data["pages"]
	fmt.Println(pages)

	items := data["items"].([]interface{})

	for _, item := range items {
		// var info map[string]interface{}
		info := item.(map[string]interface{})

		name := info["name"]
		fmt.Println(name)

		area := info["area"].(map[string]interface{})
		city := area["name"]
		fmt.Println(city)

		link := info["apply_alternate_url"]
		fmt.Println(link)

		employer := info["employer"].(map[string]interface{})
		company := employer["name"]
		fmt.Println(company)

	}

	// var child map[string]interface{}
	// child = items.(map[string]interface{})
	// fmt.Println(items)
	// for _, job := range data {
	// 	job := data.(string)
	// }

	// fmt.Println(data, errParse)

}

func main() {
	// db := PrepareDB()
	// CreateTable(db)

	ExampleScrape()

	// allJobs := ExampleScrape()
	// link := `https://hh.ru/search/vacancy?text=%28junior+OR+trainee+OR+intern%29+and+%28" +
	// 		"Go+OR+Golang+OR+Python%29&only_with_salary=false&order_by=publication_time&specialization=1" +
	// 		"&area=113&enable_snippets=true&clusters=true&experience=noExperience&salary=`

	// for _, job := range allJobs {
	// 	if ExistsVacancy(db, job) {
	// 		continue
	// 	} else {
	// 		fmt.Println("Found new job posting!")
	// 		new := fmt.Sprintf("%s\n %s\n %s\n %s\n", job.Position, job.Company, job.City, job.Link)
	// 		fmt.Println(new)

	// 		InsertVacancy(db, job)
	// 	}

	// }
	// fmt.Println(isExists)

	// fmt.Println(allJobs)
	fmt.Println("Finished")
}
