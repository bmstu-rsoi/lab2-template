package core

import (
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/library/core/ports/libraries"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

type Core struct {
	libraries libraries.Client
}

func New(lg *slog.Logger, probe *readiness.Probe, library libraries.Client) (*Core, error) {
	probe.Mark("core", true)
	lg.Warn("[startup] core ready")

	return &Core{libraries: library}, nil
}
