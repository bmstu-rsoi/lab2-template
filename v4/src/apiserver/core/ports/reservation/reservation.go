package reservation

import "context"

type Config struct {
	Host string
	Port int
}

type Client interface {
	GetUserReservations(ctx context.Context, username string) ([]Reservation, error)
}
