package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"lab2/src/flight-service/internal/repository"

	"github.com/gorilla/mux"
)

type FlightHandler struct {
	FlightRepo repository.Repository
}

func (h *FlightHandler) GetAllFlightHandler(w http.ResponseWriter, r *http.Request) {
	flightRepo := h.FlightRepo
	flight, err := flightRepo.GetAllFlights()
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(flight); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FlightHandler) GetFlightHandler(w http.ResponseWriter, r *http.Request) {
	flightRepo := h.FlightRepo
	params := mux.Vars(r)
	flight, err := flightRepo.GetFlightByNumber(params["flightNumber"])
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(flight); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FlightHandler) GetAirportHandler(w http.ResponseWriter, r *http.Request) {
	flightRepo := h.FlightRepo
	params := mux.Vars(r)
	flight, err := flightRepo.GetAirportByID(params["airportId"])
	if err != nil {
		log.Printf("failed to get flghts: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(flight); err != nil {
		log.Printf("Failed to encode flight: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
