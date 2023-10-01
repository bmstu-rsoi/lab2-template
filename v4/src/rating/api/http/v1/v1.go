package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/httpvalidator"
	"github.com/migregal/bmstu-iu7-ds-lab2/rating/core/ports/ratings"
)

type Core interface {
	GetUserRating(context.Context, string) (ratings.Rating, error)
}

func InitListener(mx *echo.Echo, core Core) error {
	gr := mx.Group("/api/v1")

	a := api{core: core}

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
