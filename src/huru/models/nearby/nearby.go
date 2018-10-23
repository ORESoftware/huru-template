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
CREATE TABLE nearby (
	id SERIAL,
    me integer,
	you integer,
	contactTime Date
);
`

// CreateTable whatever
func CreateTable() {
	db := dbs.GetDatabaseConnection()
	db.Exec(schema)
}

var nearby []Nearby

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
