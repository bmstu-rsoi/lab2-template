package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LibrariesRequest struct {
	PaginatedRequest `valid:"optional"`

	City string `query:"city" valid:"required"`
}

type LibrariesResponse struct {
	PaginatedResponse

	Items []Library `json:"items"`
}

func (a *api) GetLibraries(c echo.Context, req LibrariesRequest) error {
	data, err := a.core.GetLibraries(c.Request().Context(), req.City, req.Page, req.Size)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := LibrariesResponse{Items: []Library{}}
	for _, lib := range data.Items {
		resp.Items = append(resp.Items, Library(lib))
	}
	resp.Total = data.Total
	resp.Page = req.Page
	resp.PageSize = req.Size

	return c.JSON(http.StatusOK, &resp)
}
