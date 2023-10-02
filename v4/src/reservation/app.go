package reservation

import (
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/apiutils"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/api/http"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/config"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/core"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/services/reservationdb"
)

type App struct {
	cfg *config.Config

	http *http.Server
}

func New(lg *slog.Logger, cfg *config.Config) (*App, error) {
	a := App{cfg: cfg}

	probe := readiness.New()

	reservations, err := reservationdb.New(lg.With("module", "reservation"), cfg.Reservations, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init reservations db: %w", err)
	}

	core, err := core.New(lg.With("module", "core"), probe, reservations)
	if err != nil {
		return nil, fmt.Errorf("failed to init core: %w", err)
	}

	a.http, err = http.New(lg.With("module", "http_api"), probe, core)
	if err != nil {
		return nil, fmt.Errorf("failed to init http server: %w", err)
	}

	return &a, nil
}

func (s *App) Run(lg *slog.Logger) {
	apiutils.Serve(lg,
		apiutils.NewCallable(s.cfg.HTTPAddr, s.http),
	)
}
