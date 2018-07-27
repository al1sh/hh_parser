package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	db := PrepareDB()

	user := "user2"
	cities := []int32{1, 2, 3}

	InitUser(db, user)
	// SetString(db, user, "search", "python")
	// SetString(db, user, "experience", "0")
	// SetArray(db, user, "cities", cities)

	SetElement(db, user, "search", "go")
	SetElement(db, user, "experience", int32(0))
	SetArray(db, user, "cities", cities)

	RetrieveUser(db)

	// cities := []string{"1", "2", "1624"}
	// jobs := []string{"python+junior", "(go+OR+golang)+junior", "Project+manager+AND+English"}
	// experiences := []string{"noExperience", "between1And3"}

	// allJobs := []Vacancy{}
	// ch := make(chan []Vacancy)
	// routines := 0

	// for _, job := range jobs {
	// 	for _, experience := range experiences {
	// 		for _, city := range cities {
	// 			url := "https://api.hh.ru/vacancies?text=" + job + "&area=" + city + "&experience=" + experience + "&per_page=100&specialization=1"
	// 			// fmt.Println(url)
	// 			go ExampleScrape(url, ch)
	// 			routines++
	// 		}
	// 	}
	// }
	// for i := 0; i < routines; i++ {
	// 	// fmt.Println("appending now")
	// 	allJobs = append(allJobs, <-ch...)
	// 	// fmt.Println(len(allJobs))
	// }
	// close(ch)

	// fmt.Println("moving on")
	// // allJobs := ExampleScrape()

	// for _, position := range allJobs {
	// 	// fmt.Println(position)
	// 	if ExistsVacancy(db, position, user) {
	// 		continue
	// 	} else {
	// 		fmt.Println("******\nFound new job posting!\n******")
	// 		fmt.Println(position)

	// 		InsertVacancy(db, position, user)
	// 	}

	// }

	elapsed := time.Since(start)
	fmt.Printf("Search took %s", elapsed)

}
