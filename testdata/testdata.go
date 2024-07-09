package testdata

import (
	"encoding/json"
	"time"

	"github.com/jpicht/polyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
)

type A struct {
	A               string
	polyjson.Is[AB] `polyjson:"a,omitempty"`
}

type B struct {
	B               string
	polyjson.Is[AB] `polyjson:"b,omitempty"`
}

type CommonAB struct {
	polyjson.Common[AB]
	Timestamp time.Time
}

type C struct {
	C               string
	polyjson.Is[CD] `polyjson:"c,omitempty"`
}

type D struct {
	D                       string
	polyjson.Implements[CD] `polyjson:"d,omitempty"`
}

//easyjson:json
type CommonCD struct {
	polyjson.Common[CD]
	Timestamp  time.Time
	SomeString *string
	Additional polyjson.AdditionalFieldMap[json.RawMessage]
}

type E struct {
	E               string
	polyjson.Is[EF] `polyjson:"e,omitempty"`
}

type F struct {
	F                       string
	polyjson.Implements[EF] `polyjson:"f,omitempty"`
}

type CommonEF struct {
	polyjson.Common[EF]
}

func (c *CommonEF) AdditionalField(key string, _ *jlexer.Lexer) {}
