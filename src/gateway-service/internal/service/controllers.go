package service

import (
	"fmt"

	"lab2/src/gateway-service/internal/models"
)

func CalncelTicketController(ticketServiceAddress, bonusServiceAddress, username string) error {
	return nil
}

func UserTicketsController(ticketServiceAddress, flightServiceAddress, username string) (*[]models.TicketInfo, error) {
	tickets, err := GetUserTickets(ticketServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user tickets: %s\n", err)
	}

	ticketsInfo := make([]models.TicketInfo, 0)
	for _, ticket := range *tickets {
		flight, err := GetFlight(flightServiceAddress, ticket.FlightNumber)
		if err != nil {
			return nil, fmt.Errorf("Failed to get flights: %s", err)
		}

		airportFrom, err := GetAirport(flightServiceAddress, flight.FromAirportId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get airport: %s", err)
		}

		airportTo, err := GetAirport(flightServiceAddress, flight.ToAirportId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get airport: %s", err)
		}

		ticketInfo := models.TicketInfo{
			TicketUID:    ticket.TicketUID,
			FlightNumber: ticket.FlightNumber,
			FromAirport:  fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
			ToAirport:    fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
			Date:         flight.Date,
			Price:        ticket.Price,
			Status:       ticket.Status,
		}

		ticketsInfo = append(ticketsInfo, ticketInfo)
	}

	return &ticketsInfo, nil
}

func UserInfoController(ticketServiceAddress, flightServiceAddress, bonusServiceAddress, username string) (*models.UserInfo, error) {
	ticketsInfo, err := UserTicketsController(ticketServiceAddress, flightServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user tickets: %s", err)
	}

	privilege, err := GetPrivilegeShortInfo(bonusServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get privilege info: %s", err)
	}

	userInfo := &models.UserInfo{
		TicketsInfo: ticketsInfo,
		Privilege: &models.PrivilegeShortInfo{
			Status:  privilege.Status,
			Balance: privilege.Balance,
		},
	}

	return userInfo, nil
}

func UserPrivilegeController(bonusServiceAddress, username string) (*models.PrivilegeInfo, error) {
	privilegeShortInfo, err := GetPrivilegeShortInfo(bonusServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user tickets: %s", err)
	}

	privilegeHistory, err := GetPrivilegeHistory(bonusServiceAddress, privilegeShortInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get privilege info: %s", err)
	}

	privilegeInfo := &models.PrivilegeInfo{
		Status:  privilegeShortInfo.Status,
		Balance: privilegeShortInfo.Balance,
		History: privilegeHistory,
	}

	return privilegeInfo, nil
}

func BuyTicketController(tAddr, fAddr, bAddr, username string, info *models.BuyTicketInfo) (*models.PurchaseTicketInfo, error) {
	flight, err := GetFlight(fAddr, info.FlightNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to get flights: %s", err)
	}

	airportFrom, err := GetAirport(fAddr, flight.FromAirportId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %s", err)
	}

	airportTo, err := GetAirport(fAddr, flight.ToAirportId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %s", err)
	}

	moneyPaid := flight.Price
	bonusesPaid := 0
	diff := int(float32(info.Price) * 0.1)
	optype := "FILL_IN_BALANCE"

	if info.PaidFromBalance {
		if info.Price > 0 {
			bonusesPaid = 0
		} else {
			bonusesPaid = info.Price
		}

		moneyPaid = moneyPaid - bonusesPaid
		diff = -bonusesPaid
		optype = "DEBIT_THE_ACCOUNT"
	}

	uid, err := CreateTicket(tAddr, username, info.FlightNumber, flight.Price)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ticket: %s", err)
	}

	if !info.PaidFromBalance {
		if err := CreatePrivilege(bAddr, username, diff); err != nil {
			return nil, fmt.Errorf("Failed to get privilege info: %s", err)
		}
	}

	err = CreatePrivilegeHistoryRecord(bAddr, uid, flight.Date, optype, 1, diff)
	if err != nil {
		return nil, fmt.Errorf("Failed to create bonus history record: %s", err)
	}

	purchaseInfo := models.PurchaseTicketInfo{
		TicketUID:     uid,
		FlightNumber:  info.FlightNumber,
		FromAirport:   fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
		ToAirport:     fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
		Date:          flight.Date,
		Price:         flight.Price,
		PaidByMoney:   moneyPaid,
		PaidByBonuses: bonusesPaid,
		Status:        "PAID",
		Privilege: &models.PrivilegeShortInfo{
			Balance: diff,
			Status:  "GOLD",
		},
	}

	return &purchaseInfo, nil
}
