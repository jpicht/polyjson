package polyjson

import "errors"

var (
	ErrMultipleValues = errors.New("multiple values polyjson generated struct")
	ErrNoValue        = errors.New("no value in polyjson generated struct")
)
