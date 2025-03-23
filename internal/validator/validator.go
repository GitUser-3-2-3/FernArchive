package validator

import (
	"regexp"
	"slices"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
	"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (vldtr *Validator) Valid() bool {
	return len(vldtr.Errors) == 0
}

func (vldtr *Validator) AddErrors(key, value string) {
	if _, exists := vldtr.Errors[key]; !exists {
		vldtr.Errors[key] = value
	}
}

func (vldtr *Validator) Check(ok bool, key, value string) {
	if !ok {
		vldtr.AddErrors(key, value)
	}
}

func Matches(value string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, val := range values {
		uniqueValues[val] = true
	}
	return len(values) == len(uniqueValues)
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
