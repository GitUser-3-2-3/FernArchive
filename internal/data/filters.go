package data

import "FernArchive/internal/validator"

type Filters struct {
	Page       int
	PageSize   int
	Sort       string
	SortParams []string
}

func ValidateFilters(vldtr *validator.Validator, fltr Filters) {
	vldtr.Check(fltr.Page > 0, "page", "must be greater than zero")
	vldtr.Check(fltr.Page <= 10_000_000, "page", "must be a maximum of 10 million")

	vldtr.Check(fltr.PageSize > 0, "page_size", "must be greater than zero")
	vldtr.Check(fltr.PageSize <= 100, "page_size", "must be a maximum of 100")

	vldtr.Check(validator.PermittedValue(fltr.Sort, fltr.SortParams...), "sort", "invalid sort value")
}
