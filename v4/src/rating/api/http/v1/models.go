package v1

import (
	"github.com/migregal/bmstu-iu7-ds-lab2/pkg/httpvalidator"
)

type AuthedRequest struct {
	Username string `header:"X-User-Name" valid:"required"`
}

type ValidationErrorResponse struct {
	Message string                          `json:"message"`
	Errors  []httpvalidator.ValidationError `json:"errors"`
}
