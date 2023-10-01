package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RatingRequest struct {
	AuthedRequest `valid:"optional"`
}

type RatingResponse struct {
	Stars uint64 `json:"stars" valid:"range(0|100),required"`
}


func (a *api) GetRating(c echo.Context, req RatingRequest) error {
	data, err := a.core.GetUserRating(c.Request().Context(), req.Username)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := RatingResponse{
		Stars: uint64(data.Stars),
	}

	return c.JSON(http.StatusOK, &resp)
}
