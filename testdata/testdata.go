package testdata

import (
	"encoding/json"
	"time"

	"github.com/jpicht/polyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
)

type A struct {
	A                       string
	polyjson.Implements[AB] `polyjson:"a,omitempty"`
}

type B struct {
	B                       string
	polyjson.Implements[AB] `polyjson:"b,omitempty"`
}

func (*B) UnmarshalJSON([]byte) error { return nil }

type CommonAB struct {
	polyjson.Common[AB]
	Timestamp time.Time
}

type C struct {
	C                       string
	polyjson.Implements[CD] `polyjson:"c,omitempty"`
}

func (*C) UnmarshalEasyJSON(l *jlexer.Lexer) { l.SkipRecursive() }

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
	E                       string
	polyjson.Implements[EF] `polyjson:"e,omitempty"`
}

type F struct {
	F                       string
	polyjson.Implements[EF] `polyjson:"f,omitempty"`
}

type CommonEF struct {
	polyjson.Common[EF]
}

func (c *CommonEF) AdditionalField(key string, _ *jlexer.Lexer) {}
