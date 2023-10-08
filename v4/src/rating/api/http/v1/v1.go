package v1

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/httpvalidator"
	"github.com/migregal/bmstu-iu7-ds-lab2/rating/core/ports/ratings"
)

type Core interface {
	GetUserRating(context.Context, string) (ratings.Rating, error)
	UpdateUserRating(context.Context, string, int) error
}

func InitListener(mx *echo.Echo, lg *slog.Logger, core Core) error {
	gr := mx.Group("/api/v1")

	a := api{lg: lg, core: core}

	gr.GET("/rating", WrapRequest(lg, a.GetRating))
	gr.PATCH("/rating", WrapRequest(lg, a.UpdateRating))

	return nil
}

type api struct {
	lg   *slog.Logger
	core Core
}

func WrapRequest[T any](lg *slog.Logger, handler func(echo.Context, T) error) func(echo.Context) error {
	return func(c echo.Context) error {
		binder := &echo.DefaultBinder{}

		var req T
		if err := binder.Bind(&req, c); err != nil {
			lg.Warn("failed to bind request", "error", err)
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := binder.BindQueryParams(c, &req); err != nil {
			lg.Warn("failed to bind headers", "error", err)
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := binder.BindHeaders(c, &req); err != nil {
			lg.Warn("failed to bind headers", "error", err)
			return c.String(http.StatusBadRequest, "bad request")
		}

		if err := c.Validate(req); err != nil {
			lg.Warn("failed to validate request", "error", err)
			resp := ValidationErrorResponse{
				http.StatusText(http.StatusBadRequest),
				httpvalidator.ParseErrors(err),
			}

			return c.JSON(http.StatusBadRequest, resp)
		}

		return handler(c, req)
	}
}
