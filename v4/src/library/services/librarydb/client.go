package librarydb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/migregal/bmstu-iu7-ds-lab2/library/core/ports/libraries"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

const probeKey = "librariesdb"

type DB struct {
	db *gorm.DB
}

func New(lg *slog.Logger, cfg libraries.Config, probe *readiness.Probe) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to db: %w", err)
	}

	go runMigrations(lg, db, probe, cfg.MigrationInterval, cfg.EnableTestData)

	return &DB{db: db}, nil
}

func (d *DB) GetLibraries(
	ctx context.Context, city string, page uint64, size uint64,
) (resp libraries.Libraries, err error) {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})

	stmt := tx.Offset(int((page - 1) * size)).Limit(int(size))
	if city != "" {
		stmt = stmt.Where("city = ?", city)
	}

	var libs []Library
	if err := stmt.Find(&libs).Error; err != nil {
		tx.Rollback()

		return resp, fmt.Errorf("failed to find libraries info: %w", err)
	}

	stmt = tx.Model(&Library{})
	if city != "" {
		stmt = stmt.Where("city = ?", city)
	}

	var count int64
	if err := stmt.Count(&count).Error; err != nil {
		tx.Rollback()

		return resp, fmt.Errorf("failed to count libraries: %w", err)
	}

	resp.Total = uint64(count)
	for _, lib := range libs {
		resp.Items = append(resp.Items, libraries.Library{
			ID:      lib.LibraryID.String(),
			Name:    lib.Name,
			Address: lib.Address,
			City:    lib.City,
		})
	}

	tx.Commit()

	return resp, nil
}

func (d *DB) GetLibraryBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (resp libraries.LibraryBooks, err error) {
	tx := d.db.Begin(&sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})

	var lib Library
	if err := tx.Where("library_uid = ?", libraryID).First(&lib).Error; err != nil {
		tx.Rollback()

		return resp, fmt.Errorf("failed to find library info: %w", err)
	}

	stmt := tx.Model(&LibraryBook{}).Where("fk_library_id = ?", lib.ID)
	if !showAll {
		stmt = stmt.Where("available_count > 0")
	}

	var count int64
	if err := stmt.Count(&count).Error; err != nil {
		tx.Rollback()

		return resp, fmt.Errorf("failed to count library books info: %w", err)
	}

	stmt = tx.Offset(int((page-1)*size)).Limit(int(size)).Where("fk_library_id = ?", lib.ID)
	if !showAll {
		stmt = stmt.Where("available_count > 0")
	}

	var libraryBooks []LibraryBook
	if err := stmt.Preload("BookRef").Find(&libraryBooks).Error; err != nil {
		tx.Rollback()

		return resp, fmt.Errorf("failed to select library books info: %w", err)
	}

	resp.Total = uint64(count)
	for _, book := range libraryBooks {
		resp.Items = append(resp.Items, libraries.Book{
			ID:        book.BookRef.BookID.String(),
			Name:      book.BookRef.Name,
			Author:    book.BookRef.Author,
			Genre:     book.BookRef.Genre,
			Condition: book.BookRef.Condition,
			Available: book.AvailableCount,
		})
	}

	tx.Commit()

	return resp, nil
}
