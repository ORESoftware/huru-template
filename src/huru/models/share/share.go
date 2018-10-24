package share

import (
	"encoding/json"
	"huru/dbs"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Share The person Type (more like an object)
type Share struct {
	ID         int    `json:"id"`
	Me         int    `json:"me"`
	You        int    `json:"you"`
	FieldName  string `json:"fieldName"`
	FieldValue bool   `json:"fieldValue"`
}

const schema = `
DROP TABLE share;
CREATE TABLE share (
	id SERIAL,
    me int,
    you int,
	fieldName text,
	fieldValue boolean
);
`

var (
	mtx    sync.Mutex
	shares map[string]Share
)

// Init create collection
func Init() {
	shares = make(map[string]Share)
	mtx.Lock()
	shares["1"] = Share{ID: 1, Me: 1, You: 2, FieldName: "sharePhone", FieldValue: false}
	shares["2"] = Share{ID: 2, Me: 2, You: 1, FieldName: "shareEmail", FieldValue: true}
	shares["3"] = Share{ID: 3, Me: 1, You: 2, FieldName: "sharePhone", FieldValue: true}
	shares["4"] = Share{ID: 4, Me: 2, You: 1, FieldName: "shareEmail", FieldValue: false}
	mtx.Unlock()
}

// CreateTable whatever
func CreateTable() {

	s1 := shares["1"]
	s2 := shares["2"]

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
	mtx.Lock()
	item, ok := shares[params["id"]]
	mtx.Unlock()
	if ok {
		json.NewEncoder(w).Encode(item)
	} else {
		io.WriteString(w, "null")
	}
}

// Create create a new item
func Create(w http.ResponseWriter, r *http.Request) {
	var n Share
	json.NewDecoder(r.Body).Decode(&n)
	mtx.Lock()
	shares[strconv.Itoa(n.ID)] = n
	mtx.Unlock()
	json.NewEncoder(w).Encode(&n)
}

// Delete Delete an item
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mtx.Lock()
	_, deleted := shares[params["id"]]
	delete(shares, params["id"])
	mtx.Unlock()
	json.NewEncoder(w).Encode(deleted)
}
