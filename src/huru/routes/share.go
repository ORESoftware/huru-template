package routes

import (
	"encoding/json"
	"huru/models/share"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// PersonHandler just what it says
type ShareHandler struct{}

// PeopleInjection - injects people
type ShareInjection struct {
	Share share.Map
}

// Mount just what it says
func (h ShareHandler) Mount(router *mux.Router, v ShareInjection) {
	router.HandleFunc("/share", h.makeGetMany(v)).Methods("GET")
	router.HandleFunc("/share/{id}", h.makeGetOne(v)).Methods("GET")
	router.HandleFunc("/share/{id}", h.makeCreate(v)).Methods("POST")
	router.HandleFunc("/share/{id}", h.makeDelete(v)).Methods("DELETE")
}

// MakeGetMany Display all from the people var
func (h ShareHandler) makeGetMany(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v.Share)
	}
}

// MakeGetOne Display a single data
func (h ShareHandler) makeGetOne(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		item, ok := v.Share[params["id"]]
		mtx.Unlock()
		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	}
}

// MakeCreate create a new item
func (h ShareHandler) makeCreate(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n share.Share
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.Share[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	}
}

// MakeDelete Delete an item
func (h ShareHandler) makeDelete(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, deleted := v.Share[params["id"]]
		delete(v.Share, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(deleted)
	}
}
