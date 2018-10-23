package nearby

import (
	"encoding/json"
	"huru/dbs"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Nearby whatever
type Nearby struct {
	id          int    `json:"id,omitempty"`
	me          int    `json:"me,omitempty"`
	you         int    `json:"you,omitempty"`
	contactTime string `json:"contactTime,omitempty"`
}

var schema = `
DROP TABLE nearby;
CREATE TABLE nearby (
	id SERIAL,
    me integer,
	you integer,
	contactTime Date
);
`
var nearby []Nearby

// CreateTable whatever
func CreateTable() {

	// s1 := Nearby{id: 1, me: 1, you: 2, contactTime: strconv.Itoa(time.Now())}

	s1 := Nearby{id: 1, me: 1, you: 2, contactTime: strconv.Itoa(222)}
	nearby = append(nearby, s1)

	s2 := Nearby{id: 2, me: 2, you: 1, contactTime: strconv.Itoa(223)}
	nearby = append(nearby, s2)

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx := db.MustBegin()

	// tx.MustExec("INSERT INTO share (me, you, sharePhone) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s1)
	tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s2)
	tx.Commit()

}

// GetMany Display all from the people var
func GetMany(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(nearby)
}

// GetOne Display a single data
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range nearby {
		if strconv.Itoa(item.id) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Nearby{})
}

// Create create a new item
func Create(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	var n Nearby
	json.NewDecoder(r.Body).Decode(&n)
	// n.id = params["id"]
	nearby = append(nearby, n)
	json.NewEncoder(w).Encode(nearby)
}

// Delete Delete an item
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range nearby {
		if strconv.Itoa(item.id) == params["id"] {
			nearby = append(nearby[:index], nearby[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(nearby)
	}
}
