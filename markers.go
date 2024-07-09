package polyjson

import (
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type Is[IF any] struct{}
type Implements[IF any] struct{}
type Common[IF any] struct{}

func (Is[IF]) isIs()         {}
func (Implements[IF]) isIs() {}
func (Common[IF]) isCommon() {}

type IsInterface interface{ isIs() }
type CommonInterface interface{ isCommon() }
type EasyJSONUnmarshaller = easyjson.Unmarshaler

type AdditionalFields interface {
	AdditionalField(key string, _ *jlexer.Lexer)
}
