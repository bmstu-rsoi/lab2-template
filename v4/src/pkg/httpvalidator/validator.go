package httpvalidator

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	valid "github.com/asaskevich/govalidator"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i any) error {
	result, err := valid.ValidateStruct(i)
	if !result || err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func ParseErrors(err error) []ValidationError {
	tmp := any(err)

	internal, ok := tmp.(*echo.HTTPError)
	if !ok {
		return []ValidationError{{"internal", err.Error()}}
	}

	errs := []ValidationError{}
	for _, str := range strings.Split(internal.Message.(string), ";") {
		data := strings.SplitN(str, ":", 2)

		errs = append(errs, ValidationError{strings.TrimSpace(data[0]), strings.TrimSpace(data[1])})
	}

	return errs
}
