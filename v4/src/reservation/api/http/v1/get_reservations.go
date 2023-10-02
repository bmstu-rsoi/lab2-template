package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ReservationsRequest struct {
	AuthedRequest `valid:"optional"`

	Status string `query:"status" valid:"in(RENTED|RETURNED|EXPIRED),optional"`
}

type Reservation struct {
	ID        string    `json:"reservationUid" valid:"uuidv4,required"`
	Status    string    `json:"status" valid:"in(RENTED|RETURNED|EXPIRED)"`
	Start     time.Time `json:"-"`
	End       time.Time `json:"-"`
	BookID    string    `json:"book_id"`
	LibraryID string    `json:"library_id"`
}

func (r Reservation) MarshalJSON() ([]byte, error) {
	type Alias Reservation
	return json.Marshal(&struct {
		Alias
		Start string `json:"startDate"`
		End   string `json:"tillDate"`
	}{
		Alias: (Alias)(r),
		Start: r.Start.Format(time.DateOnly),
		End:   r.End.Format(time.DateOnly),
	})
}

func (a *api) GetReservations(c echo.Context, req ReservationsRequest) error {
	data, err := a.core.GetUserReservations(c.Request().Context(), req.Username, req.Status)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := []Reservation{}
	for _, res := range data {
		resp = append(resp, Reservation(res))
	}

	return c.JSON(http.StatusOK, &resp)
}
