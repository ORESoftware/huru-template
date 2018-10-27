package person

import (
	"encoding/json"
	"errors"
	"huru/models/person"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Inter interface {
	Zoom() string
}

// PersonHandler just what it says
type PersonHandler struct{}

// PeopleInjection - injects people
type PeopleInjection struct {
	People person.Map
}

var (
	mtx sync.Mutex
)

// Mount just what it says
func (h PersonHandler) Mount(router *mux.Router, v PeopleInjection) {
	router.HandleFunc("/people", h.makeGetMany(v)).Methods("GET")
	router.HandleFunc("/people/{id}", h.makeGetOne(v)).Methods("GET")
	router.HandleFunc("/people/{id}", h.makeCreate(v)).Methods("POST")
	router.HandleFunc("/people/{id}", h.makeDelete(v)).Methods("DELETE")
	router.HandleFunc("/people/{id}", h.makeUpdateByID(v)).Methods("PUT")
}

// MakeGetMany Display all from the people var
func (h PersonHandler) makeGetMany(v PeopleInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v.People)
	}
}

// MakeGetOne Display a single data
func (h PersonHandler) makeGetOne(v PeopleInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		item, ok := v.People[params["id"]]
		mtx.Unlock()
		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	}
}

// MakeCreate just what it says
func (h PersonHandler) makeCreate(v PeopleInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n person.Person
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.People[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	}
}

// MakeDelete just what it says
func (h PersonHandler) makeDelete(v PeopleInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, isDeletable := v.People[params["id"]]
		delete(v.People, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(isDeletable)
	}
}

func (h PersonHandler) makeUpdateByID(v PeopleInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)
		var t person.Person
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		mtx.Lock()
		item, ok := v.People[params["id"]]
		mtx.Unlock()

		if !ok {
			panic(errors.New("No item to update"))
		}

		if t.Handle != "" {
			item.Handle = t.Handle
		}

		if t.Work != "" {
			item.Work = t.Work
		}

		if t.Image != "" {
			item.Image = t.Image
		}

		if t.Firstname != "" {
			item.Firstname = t.Firstname
		}

		if t.Lastname != "" {
			item.Lastname = t.Lastname
		}

		if t.Email != "" {
			item.Email = t.Email
		}

		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	}
}
