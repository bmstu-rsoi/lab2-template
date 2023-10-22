package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"lab2/src/bonus-service/internal/handlers"
	"lab2/src/bonus-service/internal/repository"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		"postgres", 5432, "postgres", "privileges", "postgres")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	bonusRepo := repository.NewMySqlRepo(db)

	handlers := &handlers.BonusHandler{
		BonusRepo: bonusRepo,
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/bonus", handlers.CreatePrivilegeHistoryHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/bonus/privilege", handlers.CreatePrivilegeHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/bonus/{username}", handlers.GetPrivilegeByUsernameHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/bonus/history/{privilegeId}", handlers.GetHistoryByIdHandler).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")

	log.Println("Server is listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
