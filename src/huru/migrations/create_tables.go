package migrations

import (
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
)

// CreateHuruTables whatever
func CreateHuruTables() {

	person.CreateTable()
	share.CreateTable()
	nearby.CreateTable()

}
