package models

type Ticket struct {
	ID           int    `json:"id"`
	TicketUID    string `json:"ticketUid"`
	Username     string `json:"username"`
	FlightNumber string `json:"flightNumber"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}

type TicketInfo struct {
	TicketUID    string `json:"ticketUid"`
	FlightNumber string `json:"flightNumber"`
	FromAirport  string `json:"fromAirport"`
	ToAirport    string `json:"toAirport"`
	Date         string `json:"date"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}

type PurchaseTicketInfo struct {
	TicketUID     string              `json:"ticketUid"`
	FlightNumber  string              `json:"flightNumber"`
	FromAirport   string              `json:"fromAirport"`
	ToAirport     string              `json:"toAirport"`
	Date          string              `json:"date"`
	Price         int                 `json:"price"`
	PaidByMoney   int                 `json:"paidByMoney"`
	PaidByBonuses int                 `json:"paidByBonuses"`
	Status        string              `json:"status"`
	Privilege     *PrivilegeShortInfo `json:"privilege"`
}

type BuyTicketInfo struct {
	FlightNumber    string `json:"flightNumber"`
	Price           int    `json:"price"`
	PaidFromBalance bool   `json:"paidFromBalance"`
}

type Flight struct {
	ID            int    `json:"id"`
	FlightNumber  string `json:"flightNumber"`
	Date          string `json:"date"`
	FromAirportId int    `json:"fromAirportId"`
	ToAirportId   int    `json:"toAirportId"`
	Price         int    `json:"price"`
}

type FlightInfo struct {
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
	FromAirport  string `json:"fromAirport"`
	ToAirport    string `json:"toAirport"`
	Price        int    `json:"price"`
}

type FlightsLimited struct {
	Page          int           `json:"page"`
	PageSize      int           `json:"pageSize"`
	TotalElements int           `json:"totalElements"`
	Items         *[]FlightInfo `json:"items"`
}

type Airport struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type Privilege struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status"`
	Balance  int    `json:"balance"`
}

type PrivilegeHistory struct {
	ID            int    `json:"id"`
	PrivilegeID   int    `json:"privilegeId"`
	TicketUID     string `json:"ticketUid"`
	Date          string `json:"date"`
	BalanceDiff   int    `json:"balanceDiff"`
	OperationType string `json:"operationType"`
}

type PrivilegeShortInfo struct {
	Status  string `json:"status"`
	Balance int    `json:"balance"`
}

type PrivilegeInfo struct {
	Balance int                 `json:"balance"`
	Status  string              `json:"status"`
	History *[]PrivilegeHistory `json:"history"`
}

type UserInfo struct {
	Privilege   *PrivilegeShortInfo `json:"privilege"`
	TicketsInfo *[]TicketInfo       `json:"tickets"`
}
