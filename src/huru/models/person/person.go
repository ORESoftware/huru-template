package person

import (
	"huru/dbs"
	"sync"
)

// Person The person Type (more like an object)
type Person struct {
	ID        int    `json:"id,omitempty"`
	Handle    string `json:"handle,omitempty"`
	Work      string `json:"work,omitempty"`
	Image     string `json:"image,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
}

const schema = `
DROP TABLE person;
DROP INDEX IF EXISTS person_handle;
DROP INDEX IF EXISTS person_email;

CREATE TABLE person (
	id SERIAL,
	handle text,
	firstname text,
    lastname text,
    email text,
	work text,
	image text,
	personalEmail text,
	businessEmail text,
	facebook text,
	instagram text
);

CREATE UNIQUE INDEX person_handle ON person (handle);
CREATE UNIQUE INDEX person_email ON person (email);
`

// Map muh person map duh
type Map map[string]Person

var (
	mtx    sync.Mutex
	people Map
)

// Init create collection
func Init() Map {
	people = make(Map)
	mtx.Lock()
	people["1"] = Person{ID: 1, Firstname: "Alex", Lastname: "Chaz", Email: "alex@example.com"}
	people["2"] = Person{ID: 2, Firstname: "Jason", Lastname: "Statham", Email: "jason@example.com"}
	mtx.Unlock()
	return people
}

// CreateTable whatever
func CreateTable() {

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx, err := db.Begin()

	if err != nil {
		panic("could not begin transaction")
	}

	s1 := people["1"]
	s2 := people["2"]

	tx.Exec("INSERT INTO person (firstname, lastname, email) VALUES ($1, $2, $3)", s1.Firstname, s1.Lastname, s1.Email)
	tx.Exec("INSERT INTO person (firstname, lastname, email) VALUES ($1, $2, $3)", s2.Firstname, s2.Lastname, s2.Email)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	// tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:Firstname, :Lastname, :Email)", s1)
	// tx.NamedExec("INSERT INTO person (firstname, lastname, email) VALUES (:Firstname, :Lastname, :Email)", s2)
	tx.Commit()
}
