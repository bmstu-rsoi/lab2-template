package library

import "context"

type Config struct {
	Host string
	Port int
}

type Client interface {
	GetBooks(context.Context, string, bool, uint64, uint64) (LibraryBooks, error)
}
