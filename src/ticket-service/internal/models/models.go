package models

type Ticket struct {
	ID           int    `json:"id"`
	TicketUID    string `json:"ticketUid"`
	Username     string `json:"username"`
	FlightNumber string `json:"flightNumber"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}
