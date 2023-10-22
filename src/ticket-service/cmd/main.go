package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"lab2/src/ticket-service/internal/handlers"
	"lab2/src/ticket-service/internal/repository"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		"postgres", 5432, "postgres", "tickets", "postgres")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ticketRepo := repository.NewMySqlRepo(db)

	handlers := &handlers.TicketHandler{
		TicketRepo: ticketRepo,
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/tickets/{username}", handlers.GetTicketsByUsernameHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/tickets", handlers.BuyTicketHandler).Methods("POST", "OPTIONS")

	port := os.Getenv("PORT")

	log.Println("Server is listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
