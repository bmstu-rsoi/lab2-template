package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ReturnBookRequest struct {
	AuthedRequest `valid:"optional"`
	ID            string    `path:"id" valid:"uuidv4,required"`
	Condition     string    `json:"condition" valid:"optional"`
	Date          time.Time `json:"date" valid:"optional"`
}

func (a *api) ReturnBook(c echo.Context, req TakeBookRequest) error {
	if false {
		resp := ErrorResponse{}

		return c.JSON(http.StatusNotFound, &resp)
	}

	return c.NoContent(http.StatusNoContent)
}
