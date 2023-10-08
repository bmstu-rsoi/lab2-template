package libraries

import (
	"context"
	"time"
)

type Config struct {
	User              string
	Password          string
	Database          string
	Host              string
	Port              int
	MigrationInterval time.Duration `mapstructure:"migration_interval"`
	EnableTestData    bool          `mapstructure:"enable_test_data"`
}

type Client interface {
	GetLibraries(ctx context.Context, city string, page uint64, size uint64) (Libraries, error)
	GetLibrariesByIDs(ctx context.Context, ids []string) (Libraries, error)
	GetLibraryBooks(ctx context.Context, libraryID string, showAll bool, page uint64, size uint64) (LibraryBooks, error)
	GetLibraryBooksByIDs(ctx context.Context, ids []string) (LibraryBooks, error)
	TakeBookFromLibrary(ctx context.Context, libraryID, bookID string) (ReservedBook, error)
}
