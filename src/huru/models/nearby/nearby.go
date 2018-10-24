package nearby

import (
	"encoding/json"
	"huru/dbs"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Nearby whatever
type Nearby struct {
	ID          int    `json:"id,omitempty"`
	Me          int    `json:"me,omitempty"`
	You         int    `json:"you,omitempty"`
	ContactTime string `json:"contactTime,omitempty"`
}

var schema = `
DROP TABLE nearby;
CREATE TABLE nearby (
	id SERIAL,
    me integer,
	you integer,
	contactTime text
);
`

func getValues(m interface{}) []interface{} {
	v := reflect.ValueOf(m)
	result := make([]interface{}, 0, v.Len())
	for _, k := range v.MapKeys() {
		result = append(result, v.MapIndex(k).Interface())
	}
	return result
}

var (
	mtx    sync.Mutex
	nearby map[string]Nearby
)

// CreateTable whatever
func CreateTable() {

	nearby = make(map[string]Nearby)

	// s1 := Nearby{id: 1, me: 1, you: 2, contactTime: strconv.Itoa(time.Now())}

	s1 := Nearby{ID: 1, Me: 1, You: 2, ContactTime: strconv.Itoa(222)}
	nearby["1"] = s1

	s2 := Nearby{ID: 2, Me: 2, You: 1, ContactTime: strconv.Itoa(223)}
	nearby["2"] = s2

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx := db.MustBegin()

	tx.MustExec("INSERT INTO nearby (me, you, contactTime) VALUES ($1, $2, $3)", s1.Me, s1.You, s1.ContactTime)
	tx.MustExec("INSERT INTO nearby (me, you, contactTime) VALUES ($1, $2, $3)", s2.Me, s2.You, s2.ContactTime)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person

	// tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s1)
	// tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s2)
	tx.Commit()

}

// GetMany Display all from the people var
func GetMany(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(getValues(nearby))
}

// GetOne Display a single data
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item, ok := nearby[params["id"]]
	if ok {
		json.NewEncoder(w).Encode(item)
	} else {
		io.WriteString(w, "null")
	}
}

// Create create a new item
func Create(w http.ResponseWriter, r *http.Request) {
	var n Nearby
	json.NewDecoder(r.Body).Decode(&n)
	nearby[strconv.Itoa(n.ID)] = n
	json.NewEncoder(w).Encode(&n)
}

// Delete Delete an item
func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mtx.Lock()
	_, deleted := nearby[params["id"]]
	delete(nearby, params["id"])
	mtx.Lock()
	json.NewEncoder(w).Encode(deleted)
}
