package v1

import (
	"strings"
	"time"

	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/httpvalidator"
)

type AuthedRequest struct {
	Username string `header:"X-User-Name" valid:"required"`
}

type PaginatedRequest struct {
	Page uint64 `query:"page" valid:"positive_uint,optional"`
	Size uint64 `query:"size" valid:"range(0|100),optional"`
}

type PaginatedResponse struct {
	Page     uint64 `json:"page"`
	PageSize uint64 `json:"pageSize"`
	Total    uint64 `json:"totalElements"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Library struct {
	ID      string `json:"libraryUid"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}

type Book struct {
	ID        string `json:"bookUid"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	Genre     string `json:"genre"`
	Condition string `json:"condition"`
	Available uint64 `json:"availableCount"`
}

type Rating struct {
	Stars uint64 `json:"stars"`
}

type ValidationErrorResponse struct {
	Message string                          `json:"message"`
	Errors  []httpvalidator.ValidationError `json:"errors"`
}

type Time struct {
	time.Time `valid:"required"`
}

func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(time.DateOnly, s)
	return
}
