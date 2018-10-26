package routes

import (
	"encoding/json"
	"huru/dbs"
	"huru/models/person"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func getTableCreationCommands(v string) string {
	return `
		CREATE TABLE share_` + v + ` PARTITION OF share FOR VALUES IN (` + v + `);
		CREATE TABLE nearby_` + v + ` PARTITION OF nearby FOR VALUES IN (` + v + `);
	`
}

var (
	mtx sync.Mutex
)

// RegisterHandler just what it says
type RegisterHandler struct{}

// Mount just what it says
func (h RegisterHandler) Mount(router *mux.Router) {
	log.Info("mounting routes here.")
}

// RegisterNewUser just what it says
func (h RegisterHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {

	db := dbs.GetDatabaseConnection()

	tx, err := db.Begin()

	if err != nil {
		panic("could not begin transaction")
	}

	p := person.Person{Handle: "foo9"}

	var id int
	err = tx.QueryRow("INSERT INTO person (handle,email,firstname,lastname) VALUES ($1, $2, $3, $4) RETURNING ID", p.Handle, p.Email, p.Firstname, p.Lastname).Scan(&id)

	if err != nil {
		panic(err)
	}

	log.Print("id:", id)
	myStr := strconv.Itoa(id)
	log.Print("myStr:", myStr)
	db.Exec(getTableCreationCommands(myStr))

	err = tx.Commit()

	json.NewEncoder(w).Encode(struct {
		ID string
	}{
		myStr,
	})
}
