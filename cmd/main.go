package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"azure.com/ecovo/trip-search-service/cmd/handler"
	"azure.com/ecovo/trip-search-service/cmd/middleware/auth"
	"azure.com/ecovo/trip-search-service/pkg/db"
	"azure.com/ecovo/trip-search-service/pkg/pubsub"
	"azure.com/ecovo/trip-search-service/pkg/pubsub/subscription"
	"azure.com/ecovo/trip-search-service/pkg/route"
	"azure.com/ecovo/trip-search-service/pkg/search"
	"azure.com/ecovo/trip-search-service/pkg/trip"
	"github.com/ably/ably-go/ably"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"googlemaps.github.io/maps"
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

	ablyClient, err := ably.NewRealtimeClient(ably.NewClientOptions(os.Getenv("ABLY_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	ablyPubSubRepository, err := subscription.NewAblyRepository(ablyClient)
	if err != nil {
		log.Fatal(err)
	}
	pubSubService := pubsub.NewService(ablyPubSubRepository)

	var tripRepository trip.Repository

	sendMocks := false
	if os.Getenv("SEND_MOCKS") != "" {
		sendMocks, err = strconv.ParseBool(os.Getenv("SEND_MOCKS"))
		if err != nil {
			log.Fatal("SEND_MOCKS env variable must be true or false")
		}
	}

	if sendMocks {
		tripRepository, err = trip.NewMockRepository()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tripRepository, err = trip.NewRestRepository(os.Getenv("TRIP_SERVICE_DOMAIN"), os.Getenv("TRIP_SERVICE_AUTH"))
		if err != nil {
			log.Fatal(err)
		}
	}

	tripUseCase := trip.NewService(tripRepository)

	mapsClient, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	routeRepository, err := route.NewGoogleMapsRepository(mapsClient)
	if err != nil {
		log.Fatal(err)
	}
	routeUseCase := route.NewService(routeRepository)

	searchRepository, err := search.NewMongoRepository(db.Searches)
	if err != nil {
		log.Fatal(err)
	}
	searchUseCase := search.NewService(searchRepository, pubSubService, tripUseCase, routeUseCase)

	r := mux.NewRouter()

	r.Handle("/search/{id}", handler.RequestID(handler.Auth(authValidator, handler.GetSearchByID(searchUseCase)))).
		Methods("GET")
	r.Handle("/search", handler.RequestID(handler.Auth(authValidator, handler.StartSearch(searchUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/search/{id}", handler.RequestID(handler.Auth(authValidator, handler.StopSearch(searchUseCase)))).
		Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
