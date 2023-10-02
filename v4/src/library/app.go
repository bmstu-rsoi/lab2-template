package library

import (
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/library/api/http"
	"github.com/migregal/bmstu-iu7-ds-lab2/library/config"
	"github.com/migregal/bmstu-iu7-ds-lab2/library/core"
	"github.com/migregal/bmstu-iu7-ds-lab2/library/services/librarydb"
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

	libraries, err := librarydb.New(lg.With("module", "library"), cfg.Libraries, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init libraries db: %w", err)
	}

	core, err := core.New(lg.With("module", "core"), probe, libraries)
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
