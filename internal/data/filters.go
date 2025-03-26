package data

import (
	"strings"

	"FernArchive/internal/validator"
)

type Filters struct {
	Page       int
	PageSize   int
	Sort       string
	SortParams []string
}

func (fltr *Filters) sortParam() string {
	for _, param := range fltr.SortParams {
		if fltr.Sort == param {
			return strings.TrimPrefix(fltr.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + fltr.Sort)
}

func (fltr *Filters) sortOrder() string {
	if strings.HasPrefix(fltr.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (fltr *Filters) limit() int {
	return fltr.PageSize
}

func (fltr *Filters) offset() int {
	return (fltr.Page - 1) * fltr.PageSize
}

func ValidateFilters(vldtr *validator.Validator, fltr Filters) {
	vldtr.Check(fltr.Page > 0, "page", "must be greater than zero")
	vldtr.Check(fltr.Page <= 10_000_000, "page", "must be a maximum of 10 million")

	vldtr.Check(fltr.PageSize > 0, "page_size", "must be greater than zero")
	vldtr.Check(fltr.PageSize <= 100, "page_size", "must be a maximum of 100")

	vldtr.Check(validator.PermittedValue(fltr.Sort, fltr.SortParams...), "sort", "invalid sort value")
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
