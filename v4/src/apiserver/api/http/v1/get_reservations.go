package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ReservationsRequest struct {
	AuthedRequest `valid:"optional"`
}

type ReservationsResponse struct {
	ID      string    `json:"reservationUid" valid:"uuidv4,required"`
	Status  string    `json:"status," valid:"in(RENTED,RETURNED,EXPIRED)"`
	Start   time.Time `json:"-"`
	End     time.Time `json:"-"`
	Book    Book      `json:"book"`
	Library Library   `json:"library"`
}

func (r ReservationsResponse) MarshalJSON() ([]byte, error) {
	type Alias ReservationsResponse
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
	data, err := a.core.GetUserReservations(c.Request().Context(), req.Username)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &data)
}
