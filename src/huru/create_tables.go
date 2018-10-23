package main

import (
	"huru/models/person"

	"github.com/jmoiron/sqlx"
)

// Create whatever
func Create(db sqlx.DB) {

	person.CreateTable()

}
