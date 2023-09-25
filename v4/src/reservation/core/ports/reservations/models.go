package reservations

import "time"

type Reservation struct {
	ID        string
	Status    string
	Start     time.Time
	End       time.Time
	BookID    string
	LibraryID string
}
