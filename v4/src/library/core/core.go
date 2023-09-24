package core

import (
	"context"
	"fmt"
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

func (c *Core) GetLibraryBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (libraries.LibraryBooks, error) {
	books, err := c.libraries.GetLibraryBooks(ctx, libraryID, showAll, page, size)
	if err != nil {
		return libraries.LibraryBooks{}, fmt.Errorf("failed to get books: %w", err)
	}

	return books, nil
}
