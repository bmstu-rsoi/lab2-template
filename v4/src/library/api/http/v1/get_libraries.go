package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/migregal/bmstu-iu7-ds-lab2/library/core/ports/libraries"
)

type LibrariesRequest struct {
	PaginatedRequest `valid:"optional"`

	City string `query:"city" valid:"optional"`
	IDs  string `query:"ids" valid:"optional"`
}

type LibrariesResponse struct {
	PaginatedResponse

	Items []Library `json:"items"`
}

func (a *api) GetLibraries(c echo.Context, req LibrariesRequest) error {
	if req.City == "" && len(req.IDs) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	var data libraries.Libraries
	if req.City != "" {
		var err error
		data, err = a.core.GetLibraries(c.Request().Context(), req.City, req.Page, req.Size)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	} else {
		var ids []string
		err := json.Unmarshal([]byte(req.IDs), &ids)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		data, err = a.core.GetLibrariesByIDs(c.Request().Context(), ids)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	resp := LibrariesResponse{Items: make([]Library, 0, len(data.Items))}
	resp.Total = data.Total
	for _, lib := range data.Items {
		resp.Items = append(resp.Items, Library(lib))
	}

	return c.JSON(http.StatusOK, &resp)
}
