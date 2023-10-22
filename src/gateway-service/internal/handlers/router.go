package handlers

import (
	"github.com/gorilla/mux"
)

type ServicesStruct struct {
	TicketServiceAddress string
	FlightServiceAddress string
	BonusServiceAddress  string
}

type GatewayService struct {
	Config ServicesStruct
}

func NewGatewayService(config *ServicesStruct) *GatewayService {
	return &GatewayService{Config: *config}
}

func Router() *mux.Router {
	servicesConfig := ServicesStruct{
		TicketServiceAddress: "http://ticket-service:8070",
		FlightServiceAddress: "http://flight-service:8060",
		BonusServiceAddress:  "http://bonus-service:8050",
	}

	router := mux.NewRouter()
	gs := NewGatewayService(&servicesConfig)

	router.HandleFunc("/api/v1/flights", gs.GetAllFlights).Methods("GET", "OPTIONS")               // OK
	router.HandleFunc("/api/v1/me", gs.GetUserInfo).Methods("GET", "OPTIONS")                      // OK
	router.HandleFunc("/api/v1/tickets", gs.GetUserTickets).Methods("GET", "OPTIONS")              // OK
	router.HandleFunc("/api/v1/tickets/{ticketUid}", gs.GetUserTicket).Methods("GET", "OPTIONS")   // OK
	router.HandleFunc("/api/v1/tickets", gs.BuyTicket).Methods("POST", "OPTIONS")                  // OK
	router.HandleFunc("/api/v1/tickets/{ticketUid}", gs.CancelTicket).Methods("DELETE", "OPTIONS") // OK
	router.HandleFunc("/api/v1/privilege", gs.GetPrivilege).Methods("GET", "OPTIONS")              // OK

	return router
}
