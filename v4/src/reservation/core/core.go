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

func (c *Core) AddReservation(
	ctx context.Context, username string, data reservations.Reservation,
) (string, error) {
	id, err := c.reservations.AddReservation(ctx, username, data)
	if err != nil {
		return "", fmt.Errorf("failed to add user reservation: %w", err)
	}

	return id, nil
}

func (c *Core) GetUserReservations(
	ctx context.Context, username, status string,
) ([]reservations.Reservation, error) {
	data, err := c.reservations.GetUserReservations(ctx, username, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of user reservations: %w", err)
	}

	return data, nil
}

func (c *Core) UpdateUserReservation(
	ctx context.Context, id, status string,
) error {
	err := c.reservations.UpdateUserReservation(ctx, id, status)
	if err != nil {
		return fmt.Errorf("failed to get list of user reservations: %w", err)
	}

	return nil
}
