package polyjson

import (
	"encoding/json"

	"github.com/mailru/easyjson/jlexer"
)

type AdditionalFieldMap[ValueType any] map[string]ValueType

func (a AdditionalFieldMap[ValueType]) AdditionalField(key string, in *jlexer.Lexer) {
	var value ValueType
	in.AddError(json.Unmarshal(in.Raw(), &value))
	a[key] = value
}
