package register

import (
	"encoding/json"
	"fmt"
	"huru/dbs"
	"huru/models/person"
	"huru/utils"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func getTableCreationCommands(v int) string {
	return fmt.Sprintf(`
		CREATE TABLE share_%d PARTITION OF share FOR VALUES IN (%d);
		CREATE TABLE nearby_%d PARTITION OF nearby FOR VALUES IN (%d);
	`, v, v, v, v)
}

var (
	mtx sync.Mutex
)

// Handler => RegisterHandler just what it says
type Handler struct{}

// Mount just what it says
func (h Handler) Mount(router *mux.Router, v interface{}) {
	router.HandleFunc("/register", h.makeRegisterNewUser(v)).Methods("POST")
}

// MakeRegisterNewUser just what it says
func (h Handler) makeRegisterNewUser(v interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var t struct {
			Handle string
		}
		err := decoder.Decode(&t)

		if err != nil {
			panic(err)
		}

		if t.Handle == "" {
			panic(utils.AppError{
				StatusCode: 409,
				Message:    "Missing handle property in request body.",
			})
		}

		db := dbs.GetDatabaseConnection()

		tx, err := db.Begin()

		if err != nil {
			panic(err)
		}

		p := person.Model{Handle: t.Handle}

		var id int
		err = tx.QueryRow("INSERT INTO person (handle,email,firstname,lastname) VALUES ($1, $2, $3, $4) RETURNING ID", p.Handle, p.Email, p.Firstname, p.Lastname).Scan(&id)

		if err != nil {
			panic(err)
		}

		db.Exec(getTableCreationCommands(id))
		err = tx.Commit()

		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(struct {
			ID string
		}{
			fmt.Sprintf("%d", id),
		})
	}
}
