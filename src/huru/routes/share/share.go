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

// ShareHandler just what it says
type ShareHandler struct{}

// ShareInjection - injects people
type ShareInjection struct {
	Share share.Map
}

var (
	mtx sync.Mutex
)

// Mount just what it says
func (h ShareHandler) Mount(router *mux.Router, v ShareInjection) {
	router.HandleFunc("/share", h.makeGetMany(v)).Methods("GET")
	router.HandleFunc("/share/{id}", h.makeGetOne(v)).Methods("GET")
	router.HandleFunc("/share/{id}", h.makeCreate(v)).Methods("POST")
	router.HandleFunc("/share/{id}", h.makeDelete(v)).Methods("DELETE")
	router.HandleFunc("/share/{id}", h.makeUpdateByID(v)).Methods("PUT")
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

func (h ShareHandler) makeUpdateByID(v ShareInjection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		decoder := json.NewDecoder(r.Body)
		var t share.Share
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
				panic(errors.New(utils.JoinArgs("FieldName does not match, see: ", t.FieldName, item.FieldName)))
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
