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
	GetLibraryBooks(ctx context.Context, libraryID string, showAll bool, page uint64, size uint64) (LibraryBooks, error)
}
