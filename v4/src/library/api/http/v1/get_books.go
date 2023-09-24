package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type BooksRequest struct {
	PaginatedRequest

	ShawAll   bool   `query:"showAll" valid:"optional"`
	LibraryID string `param:"libraryUid" valid:"uuidv4,required"`
}

type BooksResponse struct {
	PaginatedResponse

	Items    []Library `json:"items"`
}

func (a *api) GetLibraryBooks(c echo.Context, req BooksRequest) error {
	resp := BooksResponse{}

	return c.JSON(http.StatusOK, &resp)
}
