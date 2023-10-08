package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/library"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/rating"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/core/ports/reservation"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

var (
	ErrInsufficientRating = errors.New("insufficient rating")
	ErrNotFound           = errors.New("not found")
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
) ([]reservation.ReservationFullInfo, error) {
	resvs, err := c.reservation.GetUserReservations(ctx, username, "")
	if err != nil {
		c.lg.ErrorContext(ctx, "failed to get list of user reservations", "error", err)
		return nil, fmt.Errorf("failed to get list of user reservations: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var (
		errs      = make(chan error, 2)
		libraries library.Libraries
		books     library.LibraryBooks
	)

	go func() {
		defer wg.Done()

		ids := make([]string, 0, len(resvs))
		for _, resv := range resvs {
			ids = append(ids, resv.BookID)
		}

		var err error
		if books, err = c.library.GetBooksByIDs(ctx, ids); err != nil {
			errs <- err
			return
		}
	}()
	go func() {
		defer wg.Done()

		ids := make([]string, 0, len(resvs))
		for _, resv := range resvs {
			ids = append(ids, resv.LibraryID)
		}

		var err error
		if libraries, err = c.library.GetLibrariesByIDs(ctx, ids); err != nil {
			errs <- err
			return
		}
	}()

	wg.Wait()

	select {
	case err = <-errs:
		c.lg.ErrorContext(ctx, "failed to get list of user books", "error", err)
		return nil, fmt.Errorf("failed to get list of user books: %w", err)
	default:
	}

	data := make([]reservation.ReservationFullInfo, 0, len(resvs))
	for _, resv := range resvs {
		info := reservation.ReservationFullInfo{
			ID:       resv.ID,
			Username: username,
			Status:   resv.Status,
			Start:    resv.Start,
			End:      resv.End,
		}
		for _, library := range libraries.Items {
			if resv.LibraryID == library.ID {
				info.ReservedBook.Library = library
				break
			}
		}
		for _, book := range books.Items {
			if resv.BookID == book.ID {
				info.ReservedBook.Book = book
				break
			}
		}

		data = append(data, info)
	}

	return data, nil
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
		ID:           rsvtn.ID,
		Username:     rsvtn.Username,
		Status:       rsvtn.Status,
		Start:        rsvtn.Start,
		End:          rsvtn.End,
		ReservedBook: book,
		Rating:       rating,
	}

	return res, nil
}

func (c *Core) ReturnBook(
	ctx context.Context, username, reservationID, condition string, date time.Time,
) error {
	bookIsOK := true

	resvs, err := c.reservation.GetUserReservations(ctx, username, "RENTED")
	if err != nil {
		c.lg.Warn("failed to get reservations", "error", err)
		return fmt.Errorf("failed to get user reservations: %w", err)
	}

	var resv reservation.Reservation
	for _, r := range resvs {
		if r.ID != reservationID {
			continue
		}

		resv = r
	}

	if resv.ID == "" {
		return ErrNotFound
	}

	status := "RETURNED"
	if date.After(resv.End) {
		status, bookIsOK = "EXPIRED", false

		// TODO: decrease stars
	}

	err = c.reservation.SetUserReservationStatus(ctx, reservationID, status)
	if err != nil {
		return fmt.Errorf("failed to change reservation status: %w", err)
	}

	book, err := c.library.ReturnBook(ctx, resv.LibraryID, resv.BookID)
	if err != nil {
		return fmt.Errorf("failed to obtain book info: %w", err)
	}

	if condition != book.Condition {
		bookIsOK = false

		// TODO: decrease stars
	}

	if bookIsOK {
		// TODO: increase stars
	}

	return nil
}
