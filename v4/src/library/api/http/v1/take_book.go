package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TakeBookRequest struct {
	BookID    string `json:"bookUid" valid:"uuidv4,required"`
	LibraryID string `json:"libraryUid" valid:"uuidv4,required"`
}

type TakeBookResponse struct {
	Book    Book    `json:"book"`
	Library Library `json:"library"`
}

func (a *api) TakeBook(c echo.Context, req TakeBookRequest) error {
	data, err := a.core.TakeBook(c.Request().Context(), req.LibraryID, req.BookID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := TakeBookResponse{
		Book:    Book(data.Book),
		Library: Library(data.Library),
	}

	return c.JSON(http.StatusOK, resp)
}
