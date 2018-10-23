package main

import (
	"encoding/json"
	"huru/zoom"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// https://www.cyberciti.biz/faq/howto-add-postgresql-user-account/

// Person The person Type (more like an object)
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// Address is an address
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

// GetPeople Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetPerson Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// CreatePerson create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// DeletePerson Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

// main function to boot up everything
func main() {
	XX()
	zoom.Moo()

	messages := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {

		defer wg.Done()

		db, err := sqlx.Connect("postgres", "user=tom dbname=jerry password=myPassword sslmode=disable")
		if err != nil {
			log.Fatalln(err)
		}

		// exec the schema or fail; multi-statement Exec behavior varies between
		// database drivers;  pq will exec them all, sqlite3 won't, ymmv
		db.Exec(schema)

		router := mux.NewRouter()
		// people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
		// people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
		router.HandleFunc("/people", GetPeople).Methods("GET")
		router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
		router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
		router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
		log.Fatal(http.ListenAndServe(":8000", router))
		messages <- 1
	}()

	wg.Wait()
}
