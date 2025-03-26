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
