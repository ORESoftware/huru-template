package dbs

import (
	"database/sql"
	"sync"
)

var db *sql.DB
var once sync.Once

// GetDatabaseConnection whatever
func GetDatabaseConnection() *sql.DB {

	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", "user=tom dbname=jerry password=myPassword sslmode=disable")
		if err != nil {
			panic(err)
		}

	})

	return db
}
