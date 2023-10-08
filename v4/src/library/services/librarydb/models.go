package librarydb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Library struct {
	ID      uuid.UUID `gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	Name    string    `gorm:"size:80;not null"`
	City    string    `gorm:"size:255;not null"`
	Address string    `gorm:"size:255;not null"`
}

type Book struct {
	ID        uuid.UUID `gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	Name      string    `gorm:"size:255;not null"`
	Author    string    `gorm:"size:255"`
	Genre     string    `gorm:"size:255"`
	Condition string    `gorm:"size:20;check:condition in ('EXCELLENT','GOOD','BAD');default:'EXCELLENT'"`
}

type LibraryBook struct {
	gorm.Model

	FkBookID    uuid.UUID `gorm:"index:idx_member"`
	BookRef     Book      `gorm:"foreignkey:FkBookID;references:id"`
	FkLibraryID uuid.UUID `gorm:"index:idx_member"`
	LibraryRef  Library   `gorm:"foreignkey:FkLibraryID;references:id"`

	AvailableCount uint64 `gorm:"not null;check:available_count >= 0"`
}
