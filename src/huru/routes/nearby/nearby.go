package nearby

import (
	"encoding/json"
	"errors"
	"huru/models/nearby"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// NearbyHandler just what it says
type Handler struct{}

// NearbyInjection - injects nearby
type NearbyInjection struct {
	Nearby nearby.Map
}

var (
	mtx sync.Mutex
)

// Mount just what it says
func (h Handler) Mount(router *mux.Router, v NearbyInjection) {
	router.HandleFunc("/nearby", h.makeGetMany(v)).Methods("GET")
	router.HandleFunc("/nearby/{id}", h.makeGetOne(v)).Methods("GET")
	router.HandleFunc("/nearby/{id}", h.makeCreate(v)).Methods("POST")
	router.HandleFunc("/nearby/{id}", h.makeDelete(v)).Methods("DELETE")
	router.HandleFunc("/nearby/{id}", h.makeUpdate(v)).Methods("PUT")
}

// GetMany Display all from the people var
func (h Handler) makeGetMany(v NearbyInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v.Nearby)
	}
}

// GetOne Display a single data
func (h Handler) makeGetOne(v NearbyInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		item, ok := v.Nearby[params["id"]]
		mtx.Unlock()
		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	}
}

// GetOne Display a single data
func (h Handler) makeUpdate(v NearbyInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)
		var t nearby.Model
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		mtx.Lock()
		item, ok := v.Nearby[params["id"]]
		mtx.Unlock()

		if !ok {
			panic(errors.New("No item to update"))
		}

		if t.ContactTime != 0 {
			item.ContactTime = t.ContactTime
		}

		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	}
}

// Create create a new item
func (h Handler) makeCreate(v NearbyInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n nearby.Model
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.Nearby[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	}
}

// Delete Delete an item
func (h Handler) makeDelete(v NearbyInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, deleted := v.Nearby[params["id"]]
		delete(v.Nearby, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(deleted)
	}
}
