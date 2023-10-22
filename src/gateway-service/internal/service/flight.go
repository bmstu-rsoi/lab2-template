package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"lab2/src/gateway-service/internal/models"
)

func GetFlight(flightServiceAddress, flightNumber string) (*models.Flight, error) {
	requestURL := fmt.Sprintf("%s/api/v1/flight/%s", flightServiceAddress, flightNumber)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to flight service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	flight := &models.Flight{}
	if err = json.NewDecoder(res.Body).Decode(flight); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return flight, nil
}

func GetAllFlightsInfo(flightServiceAddress string) (*[]models.FlightInfo, error) {
	requestURL := fmt.Sprintf("%s/api/v1/flights", flightServiceAddress)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to flight service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	flights := &[]models.Flight{}
	if err = json.NewDecoder(res.Body).Decode(flights); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	flightsInfo := make([]models.FlightInfo, 0)
	for _, flight := range *flights {
		airportFrom, err := GetAirport(flightServiceAddress, flight.FromAirportId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get airport: %s", err)
		}

		airportTo, err := GetAirport(flightServiceAddress, flight.ToAirportId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get airport: %s", err)
		}

		fInfo := models.FlightInfo{
			FlightNumber: flight.FlightNumber,
			FromAirport:  fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
			ToAirport:    fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
			Date:         flight.Date,
			Price:        flight.Price,
		}

		flightsInfo = append(flightsInfo, fInfo)
	}

	return &flightsInfo, nil
}

func GetAirport(flightServiceAddress string, airportID int) (*models.Airport, error) {
	requestURL := fmt.Sprintf("%s/api/v1/flight/airport/%d", flightServiceAddress, airportID)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to flight service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	airport := &models.Airport{}
	if err = json.NewDecoder(res.Body).Decode(airport); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return airport, nil
}
