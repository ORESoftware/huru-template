package routes

import (
	"encoding/json"
	"huru/dbs"
	"huru/models/person"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Share The person Type (more like an object)
type Share struct {
	ID         int    `json:"id"`
	Me         int    `json:"me"`
	You        int    `json:"you"`
	FieldName  string `json:"fieldName"`
	FieldValue bool   `json:"fieldValue"`
}

const shareTable = `
CREATE TABLE share_0 PARTITION OF share FOR VALUES IN (0);
`

func getTableCreationCommands(v string) string {
	return `
	CREATE TABLE share_` + v + ` PARTITION OF share FOR VALUES IN (` + v + `);
	CREATE TABLE nearby_` + v + ` PARTITION OF nearby FOR VALUES IN (` + v + `);
	`
}

var (
	mtx    sync.Mutex
	shares map[string]Share
)

// RegisterNewUser just what it says
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {

	db := dbs.GetDatabaseConnection()
	// defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		panic("could not begin transaction")
	}

	p := person.Person{Handle: "foo"}

	result, err := tx.Exec("INSERT INTO person (handle,email,firstname,lastname) VALUES ($1, $2, $3) RETURNING id", p.Handle, p.Email, p.Firstname, p.Lastname)

	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	db.Exec(getTableCreationCommands(strconv.Itoa(id)))

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s1)
	// tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s2)
	tx.Commit()

	json.NewEncoder(w).Encode(shares)
}
