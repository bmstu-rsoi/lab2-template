package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/library/core/ports/libraries"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/httpvalidator"
)

type Core interface {
	GetLibraryBooks(context.Context, string, bool, uint64, uint64) (libraries.LibraryBooks, error)
	GetLibraryBooksByIDs(context.Context, []string) (libraries.LibraryBooks, error)
	GetLibraries(context.Context, string, uint64, uint64) (libraries.Libraries, error)
	GetLibrariesByIDs(context.Context, []string) (libraries.Libraries, error)
	TakeBook(context.Context, string, string) (libraries.ReservedBook, error)
	ReturnBook(context.Context, string, string) (libraries.Book, error)
}

func InitListener(mx *echo.Echo, core Core) error {
	gr := mx.Group("/api/v1")

	a := api{core: core}

	gr.GET("/libraries", WrapRequest(a.GetLibraries))
	gr.GET("/libraries/:id/books", WrapRequest(a.GetLibraryBooks))
	gr.POST("/libraries/:lib_id/books/:book_id/return", WrapRequest(a.ReturnBook))

	gr.POST("/books", WrapRequest(a.TakeBook))
	gr.GET("/books", WrapRequest(a.GetLibraryBooks))

	return nil
}

type api struct {
	core Core
}

func WrapRequest[T any](handler func(echo.Context, T) error) func(echo.Context) error {
	return func(c echo.Context) error {
		binder := &echo.DefaultBinder{}

		var req T
		if err := binder.Bind(&req, c); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := binder.BindHeaders(c, &req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := c.Validate(req); err != nil {
			resp := ValidationErrorResponse{
				http.StatusText(http.StatusBadRequest),
				httpvalidator.ParseErrors(err),
			}

			return c.JSON(http.StatusBadRequest, resp)
		}

		return handler(c, req)
	}
}
