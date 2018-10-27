package main

import (
	"encoding/json"
	"huru/migrations"
	"huru/models/nearby"
	"huru/models/person"
	"huru/models/share"
	loginRouter "huru/routes/login"
	nearbyRouter "huru/routes/nearby"
	personRouter "huru/routes/person"
	registerRouter "huru/routes/register"
	shareRouter "huru/routes/share"
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
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(struct {
					ID string
				}{
					"Oh shiz we done fucked up.",
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

	// handlers := routes.HuruRouteHandlers{}
	// injection := routes.HuruInjection{}

	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.Use(errorMiddleware)

	// register and login
	{
		login := loginRouter.LoginHandler{}
		router.HandleFunc("/login", login.Login).Methods("GET")
	}

	{
		register := registerRouter.RegisterHandler{}
		register.Mount(router, struct{}{})
	}

	{
		// people
		// router.HandleFunc("/people", person.GetMany).Methods("GET")
		// router.HandleFunc("/people/{id}", person.GetOne).Methods("GET")
		// router.HandleFunc("/people/{id}", person.Create).Methods("POST")
		// router.HandleFunc("/people/{id}", person.Delete).Methods("DELETE")
		people := person.Init()
		handler := personRouter.PersonHandler{}
		handler.Mount(router, personRouter.PeopleInjection{People: people})
	}

	{
		// nearby
		// router.HandleFunc("/nearby", nearby.GetMany).Methods("GET")
		// router.HandleFunc("/nearby/{id}", nearby.GetOne).Methods("GET")
		// router.HandleFunc("/nearby/{id}", nearby.Create).Methods("POST")
		// router.HandleFunc("/nearby/{id}", nearby.Delete).Methods("DELETE")
		nearbys := nearby.Init()
		handler := nearbyRouter.NearbyHandler{}
		handler.Mount(router, nearbyRouter.NearbyInjection{Nearby: nearbys})
	}

	{
		// share
		// router.HandleFunc("/share", share.GetMany).Methods("GET")
		// router.HandleFunc("/share/{id}", share.GetOne).Methods("GET")
		// router.HandleFunc("/share/{id}", share.Create).Methods("POST")
		// router.HandleFunc("/share/{id}", share.Delete).Methods("DELETE")
		shares := share.Init()
		handler := shareRouter.ShareHandler{}
		handler.Mount(router, shareRouter.ShareInjection{Share: shares})
	}

	log.Fatal(http.ListenAndServe(":8000", router))

}
