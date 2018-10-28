package nearby

import (
	"huru/dbs"
	"reflect"
	"sync"
	"time"
)

// Model - Nearby model whatever
type Model struct {
	ID          int   `json:"id,omitempty"`
	Me          int   `json:"me,omitempty"`
	You         int   `json:"you,omitempty"`
	ContactTime int64 `json:"contactTime,omitempty"`
}

var schema = `
	DROP TABLE nearby;

	CREATE TABLE nearby (
		id SERIAL,
		me integer,
		you integer,
		contactTime bigint
	) PARTITION BY LIST(me);

	CREATE TABLE nearby_0 PARTITION OF nearby FOR VALUES IN (0);
	CREATE TABLE nearby_1 PARTITION OF nearby FOR VALUES IN (1);
	CREATE TABLE nearby_2 PARTITION OF nearby FOR VALUES IN (2);
	CREATE TABLE nearby_3 PARTITION OF nearby FOR VALUES IN (3);
`

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(1)
}

func getValues(m interface{}) []interface{} {
	v := reflect.ValueOf(m)
	result := make([]interface{}, 0, v.Len())
	for _, k := range v.MapKeys() {
		result = append(result, v.MapIndex(k).Interface())
	}
	return result
}

// Map - id to model map
type Map map[string]Model

var (
	mtx    sync.Mutex
	nearby Map
)

// Init create collection
func Init() Map {
	nearby = make(map[string]Model)
	mtx.Lock()
	nearby["1"] = Model{ID: 1, Me: 1, You: 2, ContactTime: makeTimestamp()}
	nearby["2"] = Model{ID: 2, Me: 2, You: 1, ContactTime: makeTimestamp()}
	mtx.Unlock()
	return nearby
}

// CreateTable whatever
func CreateTable() {

	// s1 := Nearby{id: 1, me: 1, you: 2, contactTime: strconv.Itoa(time.Now())}

	db := dbs.GetDatabaseConnection()
	db.Exec(schema)

	tx, err := db.Begin()

	if err != nil {
		panic("could not begin transaction")
	}

	s1 := nearby["1"]
	s2 := nearby["2"]

	tx.Exec("INSERT INTO nearby (me, you, contactTime) VALUES ($1, $2, $3)", s1.Me, s1.You, s1.ContactTime)
	tx.Exec("INSERT INTO nearby (me, you, contactTime) VALUES ($1, $2, $3)", s2.Me, s2.You, s2.ContactTime)

	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s1)
	// tx.NamedExec("INSERT INTO nearby (me, you, contactTime) VALUES (:me, :you, :contactTime)", s2)
	tx.Commit()

}
