package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type BooksRequest struct {
	PaginatedRequest `valid:"optional"`

	ShowAll   bool   `query:"showAll" valid:"optional"`
	LibraryID string `param:"id" valid:"uuidv4,required"`
}

type BooksResponse struct {
	PaginatedResponse

	Items []Book `json:"items"`
}

func (a *api) GetLibraryBooks(c echo.Context, req BooksRequest) error {
	books, err := a.core.GetLibraryBooks(c.Request().Context(), req.LibraryID, req.ShowAll, req.Page, req.Size)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := BooksResponse{Items: []Book{}}
	for _, book := range books.Books {
		resp.Items = append(resp.Items, Book(book))
	}
	resp.Page = req.Page
	resp.PageSize = req.Size
	resp.Total = books.Total

	return c.JSON(http.StatusOK, &resp)
}
