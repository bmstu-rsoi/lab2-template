package libraries

import "github.com/migregal/bmstu-iu7-ds-lab2/pkg/collections"

type Library struct {
	ID      string
	Name    string
	Address string
	City    string
}

type Libraries collections.Countable[Library]

type Book struct {
	ID        string
	Name      string
	Author    string
	Genre     string
	Condition string
	Available uint64
}

type LibraryBooks collections.Countable[Book]

type ReservedBook struct {
	Book    Book
	Library Library
}
