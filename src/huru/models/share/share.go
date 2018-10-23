package share

import (
	"encoding/json"
	"huru/dbs"
	"net/http"

	"github.com/gorilla/mux"
)

// Person The person Type (more like an object)
type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
}

var schema = `
CREATE TABLE share (
	id SERIAL,
    first_name text,
    last_name text,
    email text
);
`

var people []Person

// CreateTable whatever
func CreateTable() {

	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe"})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe"})

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)
}

// GetMany Display all from the people var
func GetMany(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetOne Display a single data
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// Create create a new item
func Create(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// Delete Delete an item
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
