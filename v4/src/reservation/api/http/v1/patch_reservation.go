package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateReservationRequest struct {
	ID     string `param:"id" valid:"uuidv4,required"`
	Status string `query:"status" valid:"in(RENTED|RETURNED|EXPIRED),required"`
}

func (a *api) UpdateReservation(c echo.Context, req UpdateReservationRequest) error {
	err := a.core.UpdateUserReservation(c.Request().Context(), req.ID, req.Status)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
