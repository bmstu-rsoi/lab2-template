package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/core/ports/reservations"
)

type AddReservationRequest struct {
	AuthedRequest `valid:"optional"`
	Status        string    `json:"status," valid:"in(RENTED|RETURNED|EXPIRED)"`
	Start         time.Time `json:"startDate" valid:"required"`
	End           time.Time `json:"tillDate" valid:"required"`
	BookID        string    `json:"book_id" valid:"uuidv4,required"`
	LibraryID     string    `json:"library_id" valid:"uuidv4,required"`
}

type AddReservationResponse struct {
	ID string `json:"reservationUid" valid:"uuidv4,required"`
}

func (a *api) AddReservation(c echo.Context, req AddReservationRequest) error {
	data := reservations.Reservation{
		Status:    req.Status,
		Start:     req.Start,
		End:       req.End,
		BookID:    req.BookID,
		LibraryID: req.LibraryID,
	}

	id, err := a.core.AddReservation(c.Request().Context(), req.Username, data)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := AddReservationResponse{
		ID: id,
	}

	return c.JSON(http.StatusOK, &resp)
}
