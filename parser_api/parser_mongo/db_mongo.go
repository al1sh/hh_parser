package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	DBName   = "Drivers"
	UserInfo = "UserInfo"
)

type Vacancy struct {
	Position string
	Company  string
	Link     string
	City     string
	Details  string
}

// type User struct {
// 	user       string  `bson:"user"`
// 	active     bool    `bson:"active"`
// 	search     string  `bson:"search"`
// 	experience int32   `bson:"experience"`
// 	cities     []int32 `bson:"cities"`
// }

type User struct {
	User       string
	Active     bool
	Search     string
	Experience int32
	Cities     []int32
}

func RetrieveUser(c *mongo.Client) {
	cursor, err := c.Database(DBName).Collection(UserInfo).Find(context.Background(), nil)
	defer cursor.Close(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.Background()) {
		currentUser := User{}
		err = cursor.Decode(&currentUser)
		// keys, err := currentUser.Keys(false)
		// fmt.Println(keys)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("after decode")

		fmt.Println(currentUser)
	}
}

func (v Vacancy) String() string {
	return fmt.Sprintf("Position: %s\n Company: %s\n City: %s\n Description: %s\n Link: %s\n", v.Position, v.Company, v.City, v.Details, v.Link)
}

// MarshalBSON does custom marshalling
func (v Vacancy) MarshalBSON() (*bson.Document, error) {
	el := bson.NewDocument(
		bson.EC.String("position", v.Position),
		bson.EC.String("company", v.Company),
		bson.EC.String("link", v.Link),
		bson.EC.String("city", v.City),
		bson.EC.String("details", v.Details),
	)

	if el == nil {
		return nil, errors.New("Could not create bson element")
	}
	return el, nil
}

func InitUser(c *mongo.Client, user string) bool {
	query := bson.NewDocument(
		bson.EC.String("user", user),
		// bson.EC.Boolean("active", true),
	)

	search, err := c.Database(DBName).Collection(UserInfo).Find(context.Background(), query)
	defer search.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if search.Next(nil) {
		return false
	}

	query = bson.NewDocument(
		bson.EC.String("user", user),
		bson.EC.Boolean("active", false),
		bson.EC.String("search", ""),
		bson.EC.Int32("experience", 0),
		bson.EC.ArrayFromElements("cities", bson.VC.Int32(1)),
	)

	_, errInsert := c.Database(DBName).Collection(UserInfo).InsertOne(context.Background(), query)
	if errInsert != nil {
		log.Fatal(errInsert)
	}
	return true
}

func SetString(c *mongo.Client, user string, field string, value string) {
	query := bson.NewDocument(
		bson.EC.SubDocumentFromElements(
			"$set",
			bson.EC.String(field, value),
		),
	)

	_, err := c.Database(DBName).Collection(UserInfo).UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user", user),
		),
		query,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func SetBool(c *mongo.Client, user string, field string, value bool) {
	query := bson.NewDocument(
		bson.EC.SubDocumentFromElements(
			"$set",
			bson.EC.Boolean(field, value),
		),
	)

	_, err := c.Database(DBName).Collection(UserInfo).UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user", user),
		),
		query,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func SetElement(c *mongo.Client, user string, field string, value interface{}) {
	query := bson.NewDocument()

	switch t := value.(type) {
	case int32:
		element := bson.EC.Int32(field, t)
		query = bson.NewDocument(
			bson.EC.SubDocumentFromElements(
				"$set",
				element,
			),
		)
	case string:
		element := bson.EC.String(field, t)
		query = bson.NewDocument(
			bson.EC.SubDocumentFromElements(
				"$set",
				element,
			),
		)
	case bool:
		element := bson.EC.Boolean(field, t)
		query = bson.NewDocument(
			bson.EC.SubDocumentFromElements(
				"$set",
				element,
			),
		)
	default:
		log.Fatal("Unsupported type passed type passed")
	}

	_, err := c.Database(DBName).Collection(UserInfo).UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user", user),
		),
		query,
	)

	if err != nil {
		log.Fatal(err)
	}
}

func SetArray(c *mongo.Client, user string, field string, value []int32) {

	bsonArray := bson.NewArray()
	for _, i := range value {
		bsonArray.Append(bson.VC.Int32(int32(i)))
	}

	// query := bson.NewDocument(
	// 	bson.EC.SubDocumentFromElements(
	// 		"$push",
	// 		bson.EC.Array(
	// 			"$each",
	// 			bsonArray,
	// 		),
	// 	),
	// )
	unsetArray := bson.NewDocument(
		bson.EC.SubDocumentFromElements(
			"$unset",
			bson.EC.String(
				"cities",
				"",
			),
		),
	)

	_, err := c.Database(DBName).Collection(UserInfo).UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user", user),
		),
		unsetArray,
	)

	if err != nil {
		log.Fatal(err)
	}

	query := bson.NewDocument(
		bson.EC.SubDocumentFromElements(
			"$push",
			bson.EC.SubDocumentFromElements(
				"cities",
				bson.EC.Array(
					"$each",
					bsonArray,
				),
			),
		),
	)

	_, err = c.Database(DBName).Collection(UserInfo).UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user", user),
		),
		query,
	)

	if err != nil {
		log.Fatal(err)
	}
}

// UpdateVacancy stores found vacancis in the db
func InsertVacancy(c *mongo.Client, v Vacancy, user string) {
	_, errInsert := c.Database(DBName).Collection(user).InsertOne(context.Background(), v)
	if errInsert != nil {
		log.Fatal(errInsert)
	}
}

// ExistsVacancy checks if given job exists in the database
func ExistsVacancy(c *mongo.Client, v Vacancy, user string) bool {

	result := c.Database(DBName).Collection(user).FindOne(context.Background(), v)

	found := Vacancy{}
	err := result.Decode(&found)
	// fmt.Println(found)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false
		}

		log.Fatal(err)
	}

	return true

}

// PrepareDB and return it
func PrepareDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// ExampleScrape scrapes given URL
func ExampleScrape(url string, ch chan []Vacancy) {
	// Request the HTML page.

	jobs := []Vacancy{}

	apiGet := http.Client{
		Timeout: time.Second * 6, // Maximum of 6 secs
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

	ch <- jobs
}
