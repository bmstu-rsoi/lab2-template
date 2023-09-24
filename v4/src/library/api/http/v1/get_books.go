package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type BooksRequest struct {
	PaginatedRequest `valid:"optional"`

	ShawAll   bool   `query:"showAll" valid:"optional"`
	LibraryID string `param:"id" valid:"uuidv4,required"`
}

type BooksResponse struct {
	PaginatedResponse

	Items []Book `json:"items"`
}

func (a *api) GetLibraryBooks(c echo.Context, req BooksRequest) error {
	data, err := a.core.GetLibraryBooks(c.Request().Context(), req.LibraryID, req.ShawAll, req.Page, req.Size)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := BooksResponse{Items: []Book{}}
	resp.Total = data.Total
	for _, book := range data.Books {
		resp.Items = append(resp.Items, Book(book))
	}

	return c.JSON(http.StatusOK, &resp)
}
