package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/library"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/rating"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/reservation"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/httpvalidator"
)

type Core interface {
	GetLibraries(context.Context, string, uint64, uint64) (library.Libraries, error)
	GetLibraryBooks(context.Context, string, bool, uint64, uint64) (library.LibraryBooks, error)
	GetUserRating(ctx context.Context, username string) (rating.Rating, error)
	GetUserReservations(context.Context, string) ([]reservation.ReservationFullInfo, error)
	TakeBook(ctx context.Context, usename, libraryID, bookID string, end time.Time) (reservation.ReservationFullInfo, error)
	ReturnBook(ctx context.Context, username, reservationID, condition string, date time.Time) error
}

func InitListener(mx *echo.Echo, core Core) error {
	gr := mx.Group("/api/v1")

	a := api{core: core}

	gr.GET("/libraries", WrapRequest(a.GetLibraries))
	gr.GET("/libraries/:id/books", WrapRequest(a.GetLibraryBooks))

	gr.GET("/reservations", WrapRequest(a.GetReservations))
	gr.POST("/reservations", WrapRequest(a.TakeBook))
	gr.POST("/reservations/:id/return", WrapRequest(a.ReturnBook))

	gr.GET("/rating", WrapRequest(a.GetRating))

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
