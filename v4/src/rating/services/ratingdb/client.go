package ratingdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
	"github.com/migregal/bmstu-iu7-ds-lab2/rating/core/ports/ratings"
)

const probeKey = "ratingsdb"

type DB struct {
	db *gorm.DB
}

func New(lg *slog.Logger, cfg ratings.Config, probe *readiness.Probe) (*DB, error) {
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

func (d *DB) GetUserRating(
	ctx context.Context, username string,
) (ratings.Rating, error) {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})

	var data Rating
	stmt := tx.Where("username = ?", username)
	if err := stmt.First(&data).Error; err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ratings.ErrNotFound
		}

		return ratings.Rating{}, fmt.Errorf("failed to find rating info: %w", err)
	}

	tx.Commit()

	resp := ratings.Rating{
		Stars: data.Stars,
	}

	return resp, nil
}
