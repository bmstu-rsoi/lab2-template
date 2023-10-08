package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/rating/core/ports/ratings"
)

type Core struct {
	rating ratings.Client
}

func New(lg *slog.Logger, probe *readiness.Probe, rating ratings.Client) (*Core, error) {
	probe.Mark("core", true)
	lg.Warn("[startup] core ready")

	return &Core{rating: rating}, nil
}

func (c *Core) GetUserRating(
	ctx context.Context, username string,
) (ratings.Rating, error) {
	data, err := c.rating.GetUserRating(ctx, username)
	if err != nil {
		if errors.Is(err, ratings.ErrNotFound) {
			return ratings.Rating{Stars: 1}, nil
		}

		return ratings.Rating{}, fmt.Errorf("failed to get user rating: %w", err)
	}

	return data, nil
}

func (c *Core) UpdateUserRating(ctx context.Context, username string, diff int) error {
	err := c.rating.UpdateUserRating(ctx, username, diff)
	if err != nil {
		return fmt.Errorf("failed to get user rating: %w", err)
	}

	return nil
}
