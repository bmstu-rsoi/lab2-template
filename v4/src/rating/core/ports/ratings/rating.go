package ratings

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
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
