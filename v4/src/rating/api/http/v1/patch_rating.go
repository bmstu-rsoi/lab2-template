package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateRatingRequest struct {
	AuthedRequest `valid:"optional"`

	Diff int `query:"diff" valid:"required"`
}

func (a *api) UpdateRating(c echo.Context, req UpdateRatingRequest) error {
	err := a.core.UpdateUserRating(c.Request().Context(), req.Username, req.Diff)
	if err != nil {
		a.lg.Error("failed to update user rating", "error", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
