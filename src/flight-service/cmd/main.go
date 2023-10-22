package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"lab2/src/flight-service/internal/handlers"
	"lab2/src/flight-service/internal/repository"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		"postgres", 5432, "postgres", "flights", "postgres")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	flightRepo := repository.NewMySqlRepo(db)

	handlers := &handlers.FlightHandler{
		FlightRepo: flightRepo,
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/flights", handlers.GetAllFlightHandler).Methods("GET")
	router.HandleFunc("/api/v1/flight/{flightNumber}", handlers.GetFlightHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/flight/airport/{airportId}", handlers.GetAirportHandler).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")

	log.Println("Server is listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
