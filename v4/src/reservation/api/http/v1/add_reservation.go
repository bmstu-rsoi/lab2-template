package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AddReservationRequest struct {
	AuthedRequest `valid:"optional"`
	Status        string    `json:"status," valid:"in(RENTED,RETURNED,EXPIRED)"`
	Start         time.Time `json:"startDate" valid:"required"`
	End           time.Time `json:"tillDate" valid:"required"`
	BookID        string    `json:"book_id" valid:"uuidv4,required"`
	LibraryID     string    `json:"library_id" valid:"uuidv4,required"`
}

type AddReservationResponse struct {
	ID string `json:"reservationUid" valid:"uuidv4,required"`
}

func (a *api) AddReservation(c echo.Context, req ReservationsRequest) error {
	resp := AddReservationResponse{}

	return c.JSON(http.StatusOK, &resp)
}
