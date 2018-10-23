package person

import (
	"encoding/json"
	"huru/dbs"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Person The person Type (more like an object)
type Person struct {
	id        int    `json:"id,omitempty"`
	firstname string `json:"firstname,omitempty"`
	lastname  string `json:"lastname,omitempty"`
	email     string `json:"email,omitempty"`
}

var schema = `
DROP TABLE person;
CREATE TABLE person (
	id SERIAL,
    firstname text,
    lastname text,
    email text
);
`

var people []Person

// CreateTable whatever
func CreateTable() {

	s1 := Person{id: 1, firstname: "Alex", lastname: "Chaz", email: "alex@example.com"}
	people = append(people, s1)
	s2 := Person{id: 2, firstname: "Jason", lastname: "Statham", email: "jason@example.com"}
	people = append(people, s2)

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx := db.MustBegin()

	// tx.MustExec("INSERT INTO share (me, you, sharePhone) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:firstname, :lastname, :email)", s1)
	tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:firstname, :lastname, :email)", s2)
	tx.Commit()
}

// GetMany Display all from the people var
func GetMany(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// GetOne Display a single data
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if strconv.Itoa(item.id) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// Create create a new item
func Create(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	// person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// Delete Delete an item
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if strconv.Itoa(item.id) == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
