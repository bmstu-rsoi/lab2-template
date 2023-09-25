package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TakeBookRequest struct {
	AuthedRequest `valid:"optional"`
	BookID        string `json:"bookUid" valid:"uuidv4,required"`
	LibraryID     string `json:"libraryUid" valid:"uuidv4,required"`
	End           Time   `json:"tillDate" valid:"required"`
}

func (r TakeBookRequest) MarshalJSON() ([]byte, error) {
	type Alias TakeBookRequest
	return json.Marshal(&struct {
		Alias
		End string `json:"tillDate"`
	}{
		Alias: (Alias)(r),
		End:   r.End.Format(time.DateOnly),
	})
}

type TakeBookResponse struct {
	ID      string    `json:"reservationUid"`
	Status  string    `json:"status"`
	Start   time.Time `json:"-"`
	End     time.Time `json:"-"`
	Book    Book      `json:"book"`
	Library Library   `json:"library"`
	Rating  Rating    `json:"rating"`
}

func (r TakeBookResponse) MarshalJSON() ([]byte, error) {
	type Alias TakeBookResponse
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

func (a *api) TakeBook(c echo.Context, req TakeBookRequest) error {
	err := a.core.TakeBook(
		c.Request().Context(), req.Username, req.LibraryID, req.BookID, req.End.Time,
	)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := TakeBookResponse{}

	return c.JSON(http.StatusOK, &resp)
}
