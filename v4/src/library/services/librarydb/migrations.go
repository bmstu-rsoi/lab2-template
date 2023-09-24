package librarydb

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/readiness"
)

type migrationItem struct {
	name  string
	model any
}

func runMigrations(
	lg *slog.Logger, db *gorm.DB, probe *readiness.Probe,
	interval time.Duration, enableTestData bool,
) {
	probe.Mark(probeKey, false)

	for {
		sqlDB, err := db.DB()
		if err != nil {
			lg.Warn("[startup] failed to ping libraries db", "error", err.Error())

			continue
		}

		if err = sqlDB.Ping(); err != nil {
			lg.Warn("[startup] failed to ping libraries db", "error", err.Error())

			continue
		}

		break
	}

	models := []migrationItem{
		{"libraries", Library{}},
		{"books", Book{}},
		{"library_books", LibraryBook{}},
	}
	for !migrateModels(lg, db, models) { //nolint: revive
		time.Sleep(interval)
	}

	sync.OnceFunc(func() {
		probe.Mark(probeKey, true)
		lg.Warn("[startup] libraries db ready")
	})()

	if enableTestData {
		initTestData(lg, db)
	}
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

func initTestData(lg *slog.Logger, db *gorm.DB) {
	// data := Library{
	// 	LibraryID: uuid.MustParse("83575e12-7ce0-48ee-9931-51919ff3c9ee"),
	// 	Name:      "Библиотека имени 7 Непьющих",
	// 	City:      "Москва",
	// 	Address:   "2-я Бауманская ул., д.5, стр.1",
	// }
	data := LibraryBook{
		LibraryRef: Library{
			LibraryID: uuid.MustParse("83575e12-7ce0-48ee-9931-51919ff3c9ee"),
			Name:      "Библиотека имени 7 Непьющих",
			City:      "Москва",
			Address:   "2-я Бауманская ул., д.5, стр.1",
		},
		BookRef: Book{
			BookID:    uuid.MustParse("f7cdc58f-2caf-4b15-9727-f89dcc629b27"),
			Name:      "Краткий курс C++ в 7 томах",
			Author:    "Бьерн Страуструп",
			Genre:     "Научная фантастика",
			Condition: "EXCELLENT",
		},
		AvailableCount: 1,
	}
	if err := db.Create(&data).Error; err != nil {
		lg.Error("failed to init test data", "error", err)
	}
}
