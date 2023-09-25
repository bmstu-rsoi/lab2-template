package core

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

func (c *Core) GetLibraries(ctx context.Context, city string, page uint64, size uint64) (library.Libraries, error) {
	data, err := c.library.GetLibraries(ctx, city, page, size)
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get list of libraries", "error", err)
		return library.Libraries{}, fmt.Errorf("failed to get list of libraries: %w", err)
	}

	return data, nil
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

func (c *Core) GetUserReservations(
	ctx context.Context, username string,
) ([]reservation.Reservation, error) {
	resvs, err := c.reservation.GetUserReservations(ctx, username, "")
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get list of user reservations", "error", err)
		return nil, fmt.Errorf("failed to get list of user reservations: %w", err)
	}

	return resvs, nil
}

func (c *Core) TakeBook(
	ctx context.Context, username, libraryID, bookID string, end time.Time,
) error {
	resvs, err := c.reservation.GetUserReservations(ctx, username, "RENTED")
	if err != nil {
		return fmt.Errorf("failed to get user reservations: %w", err)
	}

	_ = len(resvs) // TODO: check me via rating system

	err = c.library.ObtainBook(ctx, libraryID, bookID)
	if err != nil {
		return fmt.Errorf("failed to obtain book from library: %w", err)
	}

	res := reservation.Reservation{
		Username:  username,
		Status:    "RENTED",
		Start:     time.Now(),
		End:       end,
		LibraryID: libraryID,
		BookID:    bookID,
	}
	_, err = c.reservation.AddUserReservation(ctx, res)
	if err != nil {
		return fmt.Errorf("failed add reservation for obtained book: %w", err)
	}

	return nil
}
