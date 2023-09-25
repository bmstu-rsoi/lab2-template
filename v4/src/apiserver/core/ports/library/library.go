package library

import "context"

type Config struct {
	Host string
	Port int
}

type Client interface {
	GetLibraries(context.Context, string, uint64, uint64) (Libraries, error)
	GetBooks(context.Context, string, bool, uint64, uint64) (LibraryBooks, error)
	ObtainBook(context.Context, string, string) error
}
