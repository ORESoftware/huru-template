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
) PARTITION BY LIST(me);

CREATE TABLE share_0 PARTITION OF share FOR VALUES IN (0);
CREATE TABLE share_1 PARTITION OF share FOR VALUES IN (1);
CREATE TABLE share_2 PARTITION OF share FOR VALUES IN (2);
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
	s3 := shares["3"]

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx, err := db.Begin()
	if err != nil {
		panic("could not begin transaction")
	}

	tx.Exec("INSERT INTO share (me, you, fieldName, fieldValue) VALUES ($1, $2, $3, $4)", s1.Me, s1.You, s1.FieldName, s1.FieldValue)
	tx.Exec("INSERT INTO share (me, you, fieldName, fieldValue) VALUES ($1, $2, $3, $4)", s2.Me, s2.You, s2.FieldName, s2.FieldValue)
	tx.Exec("INSERT INTO share (me, you, fieldName, fieldValue) VALUES ($1, $2, $3, $4)", s3.Me, s3.You, s3.FieldName, s3.FieldValue)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// tx.NamedExec("INSERT INTO share (me, you, sharePhone, shareEmail) VALUES (:me, :you, :sharePhone, :shareEmail)", s1)
	// tx.NamedExec("INSERT INTO share (me, you, sharePhone, shareEmail) VALUES (:me, :you, :sharePhone, :shareEmail)", s2)
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
