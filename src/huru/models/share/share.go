package share

import (
	"encoding/json"
	"huru/dbs"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Share The person Type (more like an object)
type Share struct {
	id         int  `json:"id,omitempty"`
	me         int  `json:"me,omitempty"`
	you        int  `json:"you,omitempty"`
	shareEmail bool `json:"shareEmail,omitempty"`
	sharePhone bool `json:"sharePhone,omitempty"`
}

var schema = `
DROP TABLE share;
CREATE TABLE share (
	id SERIAL,
    me int,
    you int,
	shareEmail boolean,
	sharePhone boolean
);
`

var shares []Share

// CreateTable whatever
func CreateTable() {

	s1 := Share{id: 1, me: 1, you: 2, sharePhone: true, shareEmail: false}
	shares = append(shares, s1)
	s2 := Share{id: 2, me: 2, you: 1, sharePhone: false, shareEmail: true}
	shares = append(shares, s2)

	db := dbs.GetDatabaseConnection()
	db.MustExec(schema)

	tx := db.MustBegin()

	// tx.MustExec("INSERT INTO share (me, you, sharePhone) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	tx.NamedExec("INSERT INTO share (me, you, sharePhone, shareEmail) VALUES (:me, :you, :sharePhone, :shareEmail)", s1)
	tx.NamedExec("INSERT INTO share (me, you, sharePhone, shareEmail) VALUES (:me, :you, :sharePhone, :shareEmail)", s2)
	tx.Commit()

}

// GetMany Display all from the people var
func GetMany(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(shares)
}

// GetOne Display a single data
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range shares {
		if strconv.Itoa(item.id) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Share{})
}

// Create create a new item
func Create(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	var s Share
	json.NewDecoder(r.Body).Decode(&s)
	// person.ID = params["id"]
	shares = append(shares, s)
	json.NewEncoder(w).Encode(s)
}

// Delete Delete an item
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range shares {
		if strconv.Itoa(item.id) == params["id"] {
			shares = append(shares[:index], shares[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(shares)
	}
}
