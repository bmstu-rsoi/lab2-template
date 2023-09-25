package reservationdb

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model

	ReservationID uuid.UUID `gorm:"column:reservation_uid;uniqueIndex;type:uuid;default:gen_random_uuid()"`
	Username      string    `gorm:"size:80;not null"`
	BookID        uuid.UUID `gorm:"column:book_uid;type:uuid;not null;<-:create"`
	LibraryID     uuid.UUID `gorm:"column:library_uid;type:uuid;not null;<-:create"`
	Status        string    `gorm:"size:20;check:status in ('RENTED', 'RETURNED', 'EXPIRED')"`
	Start         time.Time `gorm:"column:start_date;not null"`
	End           time.Time `gorm:"column:till_date;not null"`
}
