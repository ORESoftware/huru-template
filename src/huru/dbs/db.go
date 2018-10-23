package dbs

import (
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var once sync.Once

// GetDatabaseConnection whatever
func GetDatabaseConnection() *sqlx.DB {

	once.Do(func() {
		db, err := sqlx.Connect("postgres", "user=tom dbname=jerry password=myPassword sslmode=disable")
		if err != nil {
			log.Fatalln(err)
		}
		log.Print(db)
	})

	return db
}
