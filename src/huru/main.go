package main

import (
	"huru/migrations"
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// https://www.cyberciti.biz/faq/howto-add-postgresql-user-account/

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func createCollections() {
	person.Init()
	share.Init()
	nearby.Init()
}

// main function to boot up everything
func main() {

	// db := dbs.GetDatabaseConnection()

	createCollections()

	if os.Getenv("huru_env") == "db" {
		migrations.CreateHuruTables()
	}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	// people
	router.HandleFunc("/people", person.GetMany).Methods("GET")
	router.HandleFunc("/people/{id}", person.GetOne).Methods("GET")
	router.HandleFunc("/people/{id}", person.Create).Methods("POST")
	router.HandleFunc("/people/{id}", person.Delete).Methods("DELETE")

	// nearby
	router.HandleFunc("/nearby", nearby.GetMany).Methods("GET")
	router.HandleFunc("/nearby/{id}", nearby.GetOne).Methods("GET")
	router.HandleFunc("/nearby/{id}", nearby.Create).Methods("POST")
	router.HandleFunc("/nearby/{id}", nearby.Delete).Methods("DELETE")

	// share
	router.HandleFunc("/share", share.GetMany).Methods("GET")
	router.HandleFunc("/share/{id}", share.GetOne).Methods("GET")
	router.HandleFunc("/share/{id}", share.Create).Methods("POST")
	router.HandleFunc("/share/{id}", share.Delete).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}
