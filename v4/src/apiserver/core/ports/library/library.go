package library

import "context"

type Config struct {
	Host string
	Port int
}

type Client interface {
	GetLibraries(context.Context, string, uint64, uint64) (Libraries, error)
	GetLibrariesByIDs(context.Context, []string) (Libraries, error)
	GetBooks(context.Context, string, bool, uint64, uint64) (LibraryBooks, error)
	GetBooksByIDs(context.Context, []string) (LibraryBooks, error)
	ObtainBook(context.Context, string, string) (ReservedBook, error)
	ReturnBook(context.Context, string, string) (Book, error)
}
