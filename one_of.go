package polyjson

import (
	"encoding/json"

	"github.com/launchdarkly/go-jsonstream/v3/jwriter"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type OneOf[T any] struct {
	value *T
}

type OneOfEasyJSON[T easyjson.Unmarshaler] OneOf[T]

func (v *OneOf[T]) Accept(fn func(T)) bool {
	if v == nil || v.value == nil {
		return false
	}
	fn(*v.value)
	return true
}

func (v *OneOf[T]) Get() (*T, bool) {
	if v == nil || v.value == nil {
		return nil, false
	}
	return v.value, true
}

func (v *OneOf[T]) Valid() bool {
	return v != nil && v.value != nil
}

func (v OneOf[T]) MarshalJSON() ([]byte, error) {
	if v.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(v.value)
}

func (v OneOf[T]) WriteToJSONWriter(w *jwriter.Writer) {
	if v.value == nil {
		w.Null()
		return
	}
	data, err := json.Marshal(v.value)
	w.AddError(err)
	w.Raw(data)
}

func (v *OneOf[T]) UnmarshalJSON(data []byte) error {
	if v.value == nil {
		v.value = new(T)
	}
	return json.Unmarshal(data, v.value)
}

func (v *OneOf[T]) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	if v.value == nil {
		v.value = new(T)
	}
	if uej, ok := any(v.value).(interface{ UnmarshalEasyJSON(lexer *jlexer.Lexer) }); ok {
		uej.UnmarshalEasyJSON(lexer)
		return
	}
	lexer.AddError(v.UnmarshalJSON(lexer.Raw()))
}

func (v *OneOfEasyJSON[T]) UnmarshalEasyJSON(lexer *jlexer.Lexer) {
	if v.value == nil {
		v.value = new(T)
	}
	(*v.value).UnmarshalEasyJSON(lexer)
}
