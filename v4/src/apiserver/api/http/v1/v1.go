package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/api/http/validator"
)

type Core interface {
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
		var req T
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		if err := c.Validate(req); err != nil {
			resp := ValidationErrorResponse{
				http.StatusText(http.StatusBadRequest),
				validator.ParseErrors(err),
			}

			return c.JSON(http.StatusBadRequest, resp)
		}

		return handler(c, req)
	}
}
