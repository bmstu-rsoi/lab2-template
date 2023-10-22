package repository

import (
	"database/sql"
	"fmt"

	"lab2/src/ticket-service/internal/models"
)

type Repository interface {
	GetTicketsByUsername(flightNumber string) ([]*models.Ticket, error)
	CreateTicket(*models.Ticket) error
}

type TicketRepository struct {
	db *sql.DB
}

func NewMySqlRepo(db *sql.DB) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

func (r *TicketRepository) GetTicketsByUsername(username string) ([]*models.Ticket, error) {

	var tickets []*models.Ticket
	rows, err := r.db.Query(`SELECT * FROM ticket WHERE username = $1;`, username)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %s", err)
	}

	for rows.Next() {
		ticket := new(models.Ticket)
		rows.Scan(
			&ticket.ID,
			&ticket.TicketUID,
			&ticket.Username,
			&ticket.FlightNumber,
			&ticket.Price,
			&ticket.Status)

		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %s", err)
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (r *TicketRepository) CreateTicket(ticket *models.Ticket) error {

	_, err := r.db.Query(
		`INSERT INTO ticket (ticket_uid, username, flight_number, price, status) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		ticket.TicketUID,
		ticket.Username,
		ticket.FlightNumber,
		ticket.Price,
		ticket.Status,
	)

	return err
}
