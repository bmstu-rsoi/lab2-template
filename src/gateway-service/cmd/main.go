package main

import (
	"log"
	"net/http"
	"os"

	"lab2/src/gateway-service/internal/handlers"
)

func main() {
	port := os.Getenv("PORT")
	r := handlers.Router()
	log.Println("Server is listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
