package main

import (
	"huru/models/person"
	"huru/zoom"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// https://www.cyberciti.biz/faq/howto-add-postgresql-user-account/

// main function to boot up everything
func main() {

	XX()
	zoom.Moo()

	// db := dbs.GetDatabaseConnection()

	router := mux.NewRouter()
	// people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	// people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/people", person.GetMany).Methods("GET")
	router.HandleFunc("/people/{id}", person.GetOne).Methods("GET")
	router.HandleFunc("/people/{id}", person.Create).Methods("POST")
	router.HandleFunc("/people/{id}", person.Delete).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

}
