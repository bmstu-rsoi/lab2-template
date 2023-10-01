package ratings

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
}

type Client interface {
	GetUserRating(ctx context.Context, username string) (Rating, error)
}
