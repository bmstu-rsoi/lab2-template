package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RatingRequest struct {
	AuthedRequest `valid:"optional"`
}

type RatingResponse struct {
	Rating
}

func (a *api) GetRating(c echo.Context, req RatingRequest) error {
	resp := RatingResponse{}

	return c.JSON(http.StatusNotFound, &resp)
}
