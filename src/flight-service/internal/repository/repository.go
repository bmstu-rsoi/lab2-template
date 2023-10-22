package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"lab2/src/flight-service/internal/models"
)

type Repository interface {
	GetAllFlights() ([]*models.Flight, error)
	GetFlightByNumber(flightNumber string) (*models.Flight, error)
	GetAirportByID(airportID string) (*models.Airport, error)
}

type FlightRepository struct {
	db *sql.DB
}

func NewMySqlRepo(db *sql.DB) *FlightRepository {
	return &FlightRepository{
		db: db,
	}
}

func (r *FlightRepository) GetAllFlights() ([]*models.Flight, error) {

	var flights []*models.Flight
	rows, err := r.db.Query(`SELECT * FROM flight;`)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	for rows.Next() {
		f := new(models.Flight)
		if err := rows.Scan(&f.ID, &f.FlightNumber, &f.Date, &f.FromAirportId, &f.ToAirportId, &f.Price); err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}
		flights = append(flights, f)
	}

	defer rows.Close()

	return flights, nil
}

func (r *FlightRepository) GetFlightByNumber(flightNumber string) (*models.Flight, error) {
	var flight models.Flight
	row := r.db.QueryRow(`SELECT * FROM flight WHERE flight_number = $1;`, flightNumber)
	err := row.Scan(
		&flight.ID,
		&flight.FlightNumber,
		&flight.Date,
		&flight.FromAirportId,
		&flight.ToAirportId,
		&flight.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &flight, err
		}
	}

	return &flight, nil
}

func (r *FlightRepository) GetAirportByID(airportID string) (*models.Airport, error) {

	var airport models.Airport
	row := r.db.QueryRow(`SELECT * FROM airport WHERE id = $1;`, airportID)
	err := row.Scan(
		&airport.ID,
		&airport.Name,
		&airport.City,
		&airport.Country,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &airport, err
		}
	}

	return &airport, nil
}
