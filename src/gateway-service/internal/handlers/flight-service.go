package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"lab2/src/gateway-service/internal/models"
	"lab2/src/gateway-service/internal/service"
)

func (gs *GatewayService) GetAllFlights(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	flights, err := service.GetAllFlightsInfo(gs.Config.FlightServiceAddress)
	if err != nil {
		log.Printf("failed to get response from flighst service: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pageParam := params.Get("page")
	if pageParam == "" {
		log.Println("invalid query parameter: (page)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		log.Printf("unable to convert the string into int: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		log.Println("invalid query parameter (size)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		log.Printf("unable to convert the string into int:  %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	right := page * size
	if len(*flights) < right {
		right = len(*flights)
	}

	flightsStripped := (*flights)[(page-1)*size : right]
	cars := models.FlightsLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(flightsStripped),
		Items:         &flightsStripped,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cars)
	if err != nil {
		log.Printf("failed to encode response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
