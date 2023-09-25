package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/core/ports/reservations"
)

type Core struct {
	reservations reservations.Client
}

func New(lg *slog.Logger, probe *readiness.Probe, reservation reservations.Client) (*Core, error) {
	probe.Mark("core", true)
	lg.Warn("[startup] core ready")

	return &Core{reservations: reservation}, nil
}

func (c *Core) GetUserReservations(
	ctx context.Context, username string,
) ([]reservations.Reservation, error) {
	data, err := c.reservations.GetUserReservations(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of user reservations: %w", err)
	}

	return data, nil
}
