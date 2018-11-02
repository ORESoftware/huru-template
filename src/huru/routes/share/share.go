package share

import (
	"encoding/json"
	"errors"
	"huru/models/share"
	"huru/utils"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Handler - ShareHandler just what it says
type Handler struct{}

// ShareInjection - injects people
type ShareInjection struct {
	Share share.Map
}

var (
	mtx sync.Mutex
)

// Mount just what it says
func (h Handler) Mount(router *mux.Router, v ShareInjection) Handler {
	router.HandleFunc(h.makeGetMany("/api/v1/share", v)).Methods("GET")
	router.HandleFunc("/api/v1/share/{id}", h.makeGetOne(v)).Methods("GET")
	router.HandleFunc("/api/v1/share/{id}", h.makeCreate(v)).Methods("POST")
	router.HandleFunc("/api/v1/share/{id}", h.makeDelete(v)).Methods("DELETE")
	router.HandleFunc("/api/v1/share/{id}", h.makeUpdateByID(v)).Methods("PUT")
	return h
}

// MakeGetMany Display all from the people var
func (h Handler) makeGetMany(route string, v ShareInjection) (string, http.HandlerFunc) {
	return route, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(v.Share)
	}
}

// MakeGetOne Display a single data
func (h Handler) makeGetOne(v ShareInjection) http.HandlerFunc {
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
func (h Handler) makeCreate(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var n share.Model
		json.NewDecoder(r.Body).Decode(&n)
		mtx.Lock()
		v.Share[strconv.Itoa(n.ID)] = n
		mtx.Unlock()
		json.NewEncoder(w).Encode(&n)
	}
}

// MakeDelete Delete an item
func (h Handler) makeDelete(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		mtx.Lock()
		_, deleted := v.Share[params["id"]]
		delete(v.Share, params["id"])
		mtx.Unlock()
		json.NewEncoder(w).Encode(deleted)
	}
}

func (h Handler) makeUpdateByID(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)
		var t share.Model
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		mtx.Lock()
		item, ok := v.Share[params["id"]]
		mtx.Unlock()

		if !ok {
			panic(errors.New("No item to update"))
		}

		if t.FieldName != "" {
			if t.FieldName != item.FieldName {
				panic(utils.AppError{
					StatusCode: 409,
					Message:    utils.JoinArgs("FieldName does not match, see: ", t.FieldName, item.FieldName),
				})
			}
		}

		item.FieldValue = t.FieldValue

		if ok {
			json.NewEncoder(w).Encode(item)
		} else {
			io.WriteString(w, "null")
		}
	}
}
