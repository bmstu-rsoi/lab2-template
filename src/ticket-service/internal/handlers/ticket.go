package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"lab2/src/ticket-service/internal/models"
	"lab2/src/ticket-service/internal/repository"

	"github.com/gorilla/mux"
)

type TicketHandler struct {
	TicketRepo repository.Repository
}

func (h *TicketHandler) GetTicketsByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tickets, err := h.TicketRepo.GetTicketsByUsername(params["username"])
	if err != nil {
		log.Printf("Failed to get ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tickets); err != nil {
		log.Printf("Failed to encode ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TicketHandler) BuyTicketHandler(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket

	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.TicketRepo.CreateTicket(&ticket); err != nil {
		log.Printf("Failed to create ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
