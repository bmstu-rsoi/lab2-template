package apiserver

import (
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/api/http"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/config"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core"
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

	core, err := core.New(lg, probe, nil, nil, nil)
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
