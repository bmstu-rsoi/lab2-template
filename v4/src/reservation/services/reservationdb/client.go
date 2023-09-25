package reservationdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

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

func (d *DB) GetUserReservations(
	ctx context.Context, username string,
) ([]reservations.Reservation, error) {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})

	var data []reservations.Reservation
	if err := tx.Where("username = ?", username).Find(&data).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to find reservations info: %w", err)
	}

	resp := []reservations.Reservation{}
	for _, res := range data {
		resp = append(resp, reservations.Reservation{
			ID:        res.ID,
			Status:    res.Status,
			Start:     res.Start,
			End:       res.End,
			BookID:    res.BookID,
			LibraryID: res.LibraryID,
		})
	}

	tx.Commit()

	return resp, nil
}
