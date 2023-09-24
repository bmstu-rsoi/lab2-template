package core

import (
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
