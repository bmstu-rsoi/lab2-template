package reservations

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
	GetUserReservations(ctx context.Context, username, status string) ([]Reservation, error)
	AddReservation(ctx context.Context, username string, data Reservation) (string, error)
	UpdateUserReservation(ctx context.Context, id, status string) error
}
