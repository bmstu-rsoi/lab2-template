package apiserver

import (
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/api/http"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/config"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/services/library"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/services/reservation"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/apiutils"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

type App struct {
	cfg *config.Config

	http *http.Server
}

func New(lg *slog.Logger, cfg *config.Config) (*App, error) {
	a := App{cfg: cfg}

	probe := readiness.New()

	libraryapi, err := library.New(lg, cfg.Library, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init library connection: %w", err)
	}

	reservationapi, err := reservation.New(lg, cfg.Reservation, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init library connection: %w", err)
	}

	core, err := core.New(lg, probe, libraryapi, nil, reservationapi)
	if err != nil {
		return nil, fmt.Errorf("failed to init core: %w", err)
	}

	a.http, err = http.New(lg, probe, core)
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
