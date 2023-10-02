package apiserver

import (
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/api/http"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/config"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/services/library"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/services/rating"
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

	libraryapi, err := library.New(lg.With("module", "library"), cfg.Library, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init library connection: %w", err)
	}

	ratingapi, err := rating.New(lg.With("module", "rating"), cfg.Rating, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init ratings connection: %w", err)
	}

	reservationapi, err := reservation.New(lg.With("module", "reservation"), cfg.Reservation, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init reservations connection: %w", err)
	}

	core, err := core.New(lg.With("module", "core"), probe, libraryapi, ratingapi, reservationapi)
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
