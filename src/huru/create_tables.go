package main

import (
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
)

// Create whatever
func CreateHuruTables() {

	person.CreateTable()
	share.CreateTable()
	nearby.CreateTable()

}
