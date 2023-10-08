package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReturnBookRequest struct {
	LibraryID     string `path:"lib_id" valid:"uuidv4,required"`
	BookID        string `path:"book_id" valid:"uuidv4,required"`
}

type ReturnBookResponse struct {
	Book Book `path:"book"`
}

func (a *api) ReturnBook(c echo.Context, req ReturnBookRequest) error {
	data, err := a.core.ReturnBook(c.Request().Context(), req.LibraryID, req.BookID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := ReturnBookResponse{
		Book: Book(data),
	}

	return c.JSON(http.StatusOK, resp)
}
