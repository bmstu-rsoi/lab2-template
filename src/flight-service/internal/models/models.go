package models

type Flight struct {
	ID            int    `json:"id"`
	FlightNumber  string `json:"flightNumber"`
	Date          string `json:"date"`
	FromAirportId int    `json:"fromAirportId"`
	ToAirportId   int    `json:"toAirportId"`
	Price         int    `json:"price"`
}

type Airport struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
}
