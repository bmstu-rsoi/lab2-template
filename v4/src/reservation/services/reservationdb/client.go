package reservationdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/reservation/core/ports/reservations"
)

const probeKey = "reservationsdb"

type DB struct {
	db *gorm.DB
}

func New(lg *slog.Logger, cfg reservations.Config, probe *readiness.Probe) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to db: %w", err)
	}

	go runMigrations(lg, db, probe, cfg.MigrationInterval)

	return &DB{db: db}, nil
}

func (d *DB) AddReservation(
	ctx context.Context, username string, data reservations.Reservation,
) (string, error) {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable})

	bookID, err := uuid.Parse(data.BookID)
	if err != nil {
		return "", fmt.Errorf("failed to parse book id: %w", err)
	}

	libraryID, err := uuid.Parse(data.LibraryID)
	if err != nil {
		return "", fmt.Errorf("failed to parse library id: %w", err)
	}

	r := Reservation{
		Username:  username,
		BookID:    bookID,
		LibraryID: libraryID,
		Status:    data.Status,
		Start:     data.Start,
		End:       data.End,
	}
	if err := tx.Create(&r).Error; err != nil {
		tx.Rollback()

		return "", fmt.Errorf("failed to create reservation: %w", err)
	}

	tx.Commit()
	return r.ReservationID.String(), nil
}

func (d *DB) GetUserReservations(
	ctx context.Context, username, status string,
) ([]reservations.Reservation, error) {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})

	var data []Reservation
	stmt := tx.Where("username = ?", username)
	if status != "" {
		stmt = stmt.Where("status = ?", status)
	}
	if err := stmt.Find(&data).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to find reservations info: %w", err)
	}

	resp := []reservations.Reservation{}
	for _, res := range data {
		resp = append(resp, reservations.Reservation{
			ID:        res.ReservationID.String(),
			Status:    res.Status,
			Start:     res.Start,
			End:       res.End,
			BookID:    res.BookID.String(),
			LibraryID: res.LibraryID.String(),
		})
	}

	tx.Commit()

	return resp, nil
}

func (d *DB) UpdateUserReservation(ctx context.Context, id, status string) error {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelSerializable})

	stmt := tx.Model(&Reservation{}).Where("reservation_uid = ?", id)
	if err := stmt.Update("status", status).Error; err != nil {
		tx.Rollback()

		return fmt.Errorf("failed to find reservations info: %w", err)
	}

	tx.Commit()

	return nil
}
