package main

import (
	"encoding/json"
	"fmt"
	"huru/migrations"
	"huru/models"
	"huru/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				log.Error("Caught error in defer/recover middleware: ", err)
				originalError := err.(struct{ OriginalError error }).OriginalError

				if originalError != nil {
					log.Error("Original error in defer/recover middleware: ", originalError)
				}

				statusCode := err.(struct{ StatusCode int }).StatusCode

				if statusCode != 0 {
					w.WriteHeader(statusCode)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}

				message := err.(struct{ Message string }).Message

				if message == "" {
					message = "Unknown error message."
				}

				json.NewEncoder(w).Encode(struct {
					ID string
				}{
					message,
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// main function to boot up everything
func main() {

	// db := dbs.GetDatabaseConnection()

	if true || os.Getenv("huru_env") == "db" {
		migrations.CreateHuruTables()
	}

	// handlers := routes.HuruRouteHandlers{}.GetHandlers()
	// injections := routes.HuruInjection{}.GetInjections(
	// 	person.Init(),
	// 	nearby.Init(),
	// 	share.Init(),
	// )

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.Use(errorMiddleware)

	// register and login
	{
		handler := routes.LoginHandler{}
		router.HandleFunc("/login", handler.Login).Methods("GET")
	}

	{
		handler := routes.RegisterHandler{}
		handler.Mount(router, struct{}{})
	}

	{
		// people
		// router.HandleFunc("/people", person.GetMany).Methods("GET")
		// router.HandleFunc("/people/{id}", person.GetOne).Methods("GET")
		// router.HandleFunc("/people/{id}", person.Create).Methods("POST")
		// router.HandleFunc("/people/{id}", person.Delete).Methods("DELETE")

		handler := routes.PersonHandler{}
		handler.Mount(router, routes.PersonInjection{People: models.PersonInit()})
	}

	{
		// nearby
		// router.HandleFunc("/nearby", nearby.GetMany).Methods("GET")
		// router.HandleFunc("/nearby/{id}", nearby.GetOne).Methods("GET")
		// router.HandleFunc("/nearby/{id}", nearby.Create).Methods("POST")
		// router.HandleFunc("/nearby/{id}", nearby.Delete).Methods("DELETE")
		handler := routes.NearbyHandler{}
		handler.Mount(router, routes.NearbyInjection{Nearby: models.NearbyInit()})
	}

	{
		// share
		// router.HandleFunc("/share", share.GetMany).Methods("GET")
		// router.HandleFunc("/share/{id}", share.GetOne).Methods("GET")
		// router.HandleFunc("/share/{id}", share.Create).Methods("POST")
		// router.HandleFunc("/share/{id}", share.Delete).Methods("DELETE")
		handler := routes.ShareHandler{}
		handler.Mount(router, routes.ShareInjection{Share: models.ShareInit()})
	}

	host := os.Getenv("huru_api_host")
	port := os.Getenv("huru_api_port")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "8000"
	}

	log.Info(fmt.Sprintf("Huru API server listening on port %s", port))
	path := fmt.Sprintf("%s:%s", host, port)
	log.Fatal(http.ListenAndServe(path, router))

}
