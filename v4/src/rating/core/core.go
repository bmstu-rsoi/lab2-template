package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/rating/core/ports/ratings"
)

type Core struct {
	rating ratings.Client
}

func New(lg *slog.Logger, probe *readiness.Probe, reservation ratings.Client) (*Core, error) {
	probe.Mark("core", true)
	lg.Warn("[startup] core ready")

	return &Core{rating: reservation}, nil
}

func (c *Core) GetUserRating(
	ctx context.Context, username string,
) (ratings.Rating, error) {
	data, err := c.rating.GetUserRating(ctx, username)
	if err != nil {
		return ratings.Rating{}, fmt.Errorf("failed to get list of user reservations: %w", err)
	}

	return data, nil
}
