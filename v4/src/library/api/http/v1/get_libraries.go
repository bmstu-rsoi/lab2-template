package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LibrariesRequest struct {
	PaginatedRequest

	City string `query:"city" valid:"required"`
}

type LibrariesResponse struct {
	PaginatedResponse

	Items    []Library `json:"items"`
}

func (a *api) GetLibraries(c echo.Context, req LibrariesRequest) error {
	resp := LibrariesResponse{}

	return c.JSON(http.StatusOK, &resp)
}
