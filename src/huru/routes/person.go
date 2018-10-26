package routes

import (
	"encoding/json"
	"huru/models/person"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// PersonHandler just what it says
type PersonHandler struct{}

// PeopleInjection - injects people
type PeopleInjection struct {
	People person.Map
}

// Mount just what it says
func (h PersonHandler) Mount(router *mux.Router, v PeopleInjection) {
	log.Info("mounting routes here 1.")
	router.HandleFunc("/people", h.MakeGetMany(v)).Methods("GET")
	router.HandleFunc("/people/{id}", h.MakeGetOne(v)).Methods("GET")
	router.HandleFunc("/people/{id}", h.MakeCreate(v)).Methods("POST")
	router.HandleFunc("/people/{id}", h.MakeDelete(v)).Methods("DELETE")
}

// MakeGetMany Display all from the people var
func (h PersonHandler) MakeGetMany(v PeopleInjection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v.People)
	}
}

// MakeGetOne Display a single data
func (h PersonHandler) MakeGetOne(v PeopleInjection) func(http.ResponseWriter, *http.Request) {
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
func (h PersonHandler) MakeCreate(v PeopleInjection) func(http.ResponseWriter, *http.Request) {
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
func (h PersonHandler) MakeDelete(v PeopleInjection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, isDeletable := v.People[params["id"]]
		delete(v.People, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(isDeletable)
	}
}
