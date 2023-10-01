package ratingdb

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

type migrationItem struct {
	name  string
	model any
}

func runMigrations(lg *slog.Logger, db *gorm.DB, probe *readiness.Probe, interval time.Duration) {
	probe.Mark(probeKey, false)

	for {
		sqlDB, err := db.DB()
		if err != nil {
			lg.Warn("[startup] failed to ping ratings db", "error", err.Error())

			continue
		}

		if err = sqlDB.Ping(); err != nil {
			lg.Warn("[startup] failed to ping ratings db", "error", err.Error())

			continue
		}

		break
	}

	models := []migrationItem{
		{"rating", Rating{}},
	}
	for !migrateModels(lg, db, models) { //nolint: revive
		time.Sleep(interval)
	}

	sync.OnceFunc(func() {
		probe.Mark(probeKey, true)
		lg.Warn("[startup] reservations db ready")
	})()
}

func migrateModels(lg *slog.Logger, db *gorm.DB, models []migrationItem) bool {
	tx := db.Begin()

	for _, v := range models {
		v := v
		if err := db.AutoMigrate(&v.model); err != nil {
			lg.Warn(fmt.Sprintf("[startup] failed to migrate %s", v.name), "err", err)
			tx.Rollback()

			return false
		}
	}

	if err := tx.Commit().Error; err != nil {
		lg.Warn("[startup] failed to commit transaction", "err", err)

		return false
	}

	return true
}
