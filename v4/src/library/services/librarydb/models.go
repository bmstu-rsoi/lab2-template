package librarydb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Library struct {
	gorm.Model

	LibraryID uuid.UUID `gorm:"column:library_uid;uniqueIndex;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"size:80;not null"`
	City      string    `gorm:"size:255;not null"`
	Address   string    `gorm:"size:255;not null"`
}

type Book struct {
	gorm.Model

	BookID    uuid.UUID `gorm:"column:book_uid;uniqueIndex;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"size:255;not null"`
	Author    string    `gorm:"size:255"`
	Genre     string    `gorm:"size:255"`
	Condition string    `gorm:"size:20;check:condition in ('EXCELLENT','GOOD','BAD');default:'EXCELLENT'"`
}

type LibraryBook struct {
	gorm.Model

	BookID uint32
	Book
	LibraryID uint32
	Library
	AvailableCount uint32 `gorm:"not null"`
}
