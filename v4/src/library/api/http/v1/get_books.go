package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/library/core/ports/libraries"
)

type BooksRequest struct {
	PaginatedRequest `valid:"optional"`

	ShawAll   bool   `query:"showAll" valid:"optional"`
	LibraryID string `param:"id" valid:"uuidv4,optional"`
	IDs       string `query:"ids" valid:"optional"`
}

type BooksResponse struct {
	PaginatedResponse

	Items []Book `json:"items"`
}

func (a *api) GetLibraryBooks(c echo.Context, req BooksRequest) error {
	if req.LibraryID == "" && len(req.IDs) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	var data libraries.LibraryBooks
	if req.LibraryID != "" {
		var err error

		data, err = a.core.GetLibraryBooks(c.Request().Context(), req.LibraryID, req.ShawAll, req.Page, req.Size)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	} else {
		var ids []string
		err := json.Unmarshal([]byte(req.IDs), &ids)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		data, err = a.core.GetLibraryBooksByIDs(c.Request().Context(), ids)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	resp := BooksResponse{Items: make([]Book, 0, len(data.Items))}
	resp.Total = data.Total
	for _, book := range data.Items {
		resp.Items = append(resp.Items, Book(book))
	}

	return c.JSON(http.StatusOK, &resp)
}
