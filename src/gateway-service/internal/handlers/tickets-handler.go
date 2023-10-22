package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"lab2/src/gateway-service/internal/models"
	"lab2/src/gateway-service/internal/service"

	"github.com/gorilla/mux"
)

func (gs *GatewayService) GetUserTickets(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketsInfo, err := service.UserTicketsController(
		gs.Config.TicketServiceAddress,
		gs.Config.FlightServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticketsInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) CancelTicket(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := service.CalncelTicketController(
		gs.Config.TicketServiceAddress,
		gs.Config.BonusServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(204)
}

func (gs *GatewayService) GetUserTicket(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	ticketUID := params["ticketUid"]

	ticketsInfo, err := service.UserTicketsController(
		gs.Config.TicketServiceAddress,
		gs.Config.FlightServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ticketInfo *models.TicketInfo
	for _, ticket := range *ticketsInfo {
		if ticket.TicketUID == ticketUID {
			ticketInfo = &ticket
		}
	}

	if ticketInfo == nil {
		log.Printf("Ticket not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticketInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) BuyTicket(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ticketInfo models.BuyTicketInfo
	err := json.NewDecoder(r.Body).Decode(&ticketInfo)
	if err != nil {
		log.Printf("failed to decode post request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tickets, err := service.BuyTicketController(
		gs.Config.TicketServiceAddress,
		gs.Config.FlightServiceAddress,
		gs.Config.BonusServiceAddress,
		username,
		&ticketInfo,
	)

	if err != nil {
		log.Printf("failed to get response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tickets)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (gs *GatewayService) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, err := service.UserInfoController(
		gs.Config.TicketServiceAddress,
		gs.Config.FlightServiceAddress,
		gs.Config.BonusServiceAddress,
		username,
	)

	if err != nil {
		log.Printf("Failed to get response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
