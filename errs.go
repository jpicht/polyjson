package polyjson

import "errors"

var (
	ErrMultipleValues = errors.New("multiple values polyjson generated struct")
	ErrNoValue        = errors.New("no value in polyjson generated struct")
	ErrMissingTypeID  = errors.New("missing type ID field")
	ErrInvalidTypeID  = errors.New("invalid type ID value")
)
