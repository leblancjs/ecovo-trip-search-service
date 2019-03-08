package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/trip-search-service/cmd/handler"
	"azure.com/ecovo/trip-search-service/cmd/middleware/auth"
	"azure.com/ecovo/trip-search-service/pkg/db"
	"azure.com/ecovo/trip-search-service/pkg/search"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authConfig := auth.Config{
		Domain: os.Getenv("AUTH_DOMAIN")}
	authValidator, err := auth.NewTokenValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}

	dbConnectionTimeout, err := time.ParseDuration(os.Getenv("DB_CONNECTION_TIMEOUT") + "s")
	if err != nil {
		dbConnectionTimeout = db.DefaultConnectionTimeout
	}
	dbConfig := db.Config{
		Host:              os.Getenv("DB_HOST"),
		Username:          os.Getenv("DB_USERNAME"),
		Password:          os.Getenv("DB_PASSWORD"),
		Name:              os.Getenv("DB_NAME"),
		ConnectionTimeout: dbConnectionTimeout}
	db, err := db.New(&dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	searchRepository, err := search.NewMongoRepository(db.Searches)
	if err != nil {
		log.Fatal(err)
	}
	searchUseCase := search.NewService(searchRepository)

	r := mux.NewRouter()

	r.Handle("/search/{id}", handler.RequestID(handler.Auth(authValidator, handler.GetSearchByID(searchUseCase)))).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.Handle("/search", handler.RequestID(handler.Auth(authValidator, handler.StartSearch(searchUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/search/{id}", handler.RequestID(handler.Auth(authValidator, handler.StopSearch(searchUseCase)))).
		Methods("DELETE").
		Headers("Content-Type", "application/json")

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
