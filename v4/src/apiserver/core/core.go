package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/library"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/rating"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/reservation"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

type Core struct {
	lg *slog.Logger

	library     library.Client
	rating      rating.Client
	reservation reservation.Client
}

func New(
	lg *slog.Logger, probe *readiness.Probe,
	library library.Client, rating rating.Client, reservation reservation.Client,
) (*Core, error) {
	probe.Mark("core", true)
	lg.Warn("[startup] core ready")

	return &Core{lg: lg, library: library, rating: rating, reservation: reservation}, nil
}

func (c *Core) GetLibraryBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (library.LibraryBooks, error) {
	books, err := c.library.GetBooks(ctx, libraryID, showAll, page, size)
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get list of library books", "error", err)
		return library.LibraryBooks{}, fmt.Errorf("failed to get list of library books: %w", err)
	}

	return books, nil
}
