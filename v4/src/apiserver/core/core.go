package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/library"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/rating"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/reservation"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

var (
	ErrInsufficientRating = errors.New("insufficient rating")
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

func (c *Core) GetUserRating(
	ctx context.Context, username string,
) (rating.Rating, error) {
	data, err := c.rating.GetUserRating(ctx, username)
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get user rating", "error", err)
		return rating.Rating{}, fmt.Errorf("failed to get user rating: %w", err)
	}

	return data, nil
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
) (reservation.ReservationFullInfo, error) {
	resvs, err := c.reservation.GetUserReservations(ctx, username, "RENTED")
	if err != nil {
		c.lg.Warn("failed to get reservations", "error", err)
		return reservation.ReservationFullInfo{}, fmt.Errorf("failed to get user reservations: %w", err)
	}

	rating, err := c.rating.GetUserRating(ctx, username)
	if err != nil {
		c.lg.Warn("failed to get rating", "error", err)
		return reservation.ReservationFullInfo{}, fmt.Errorf("failed to get user rating: %w", err)
	}

	if uint64(len(resvs)) >= rating.Stars {
		c.lg.Warn("insufficient rating", "rating", rating.Stars)
		return reservation.ReservationFullInfo{}, ErrInsufficientRating
	}

	rsvtn := reservation.Reservation{
		Username:  username,
		Status:    "RENTED",
		Start:     time.Now(),
		End:       end,
		LibraryID: libraryID,
		BookID:    bookID,
	}

	rsvtn.ID, err = c.reservation.AddUserReservation(ctx, rsvtn)
	if err != nil {
		c.lg.Warn("failed to add reservation", "error", err)
		return reservation.ReservationFullInfo{}, fmt.Errorf("failed to add reservation for obtained book: %w", err)
	}

	book, err := c.library.ObtainBook(ctx, libraryID, bookID)
	if err != nil {
		c.lg.Warn("failed to update books amount", "error", err)
		return reservation.ReservationFullInfo{}, fmt.Errorf("failed to obtain book from library: %w", err)
	}

	res := reservation.ReservationFullInfo{
		Username:     rsvtn.Username,
		Status:       rsvtn.Status,
		Start:        rsvtn.Start,
		End:          rsvtn.End,
		ReservedBook: book,
		Rating:       rating,
	}

	return res, nil
}
