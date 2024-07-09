package polyjson

import "errors"

var (
	ErrMultipleItems  = errors.New("multiple items in OneOf struct")
	ErrMultipleValues = errors.New("multiple values in OneOf struct")
	ErrNoValue        = errors.New("no value in OneOf struct")
)
