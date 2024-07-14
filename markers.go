package polyjson

import (
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type Interface[IF any] interface{ isInterface() }
type Implements[IF any] struct{}
type Common[IF any] struct{}

func (Implements[IF]) isImplementation() {}
func (Common[IF]) isCommon()             {}

type InterfaceInterface interface{ isInterface() }
type ImplementsInterface interface{ isImplementation() }
type CommonInterface interface{ isCommon() }

type EasyJSONUnmarshaller = easyjson.Unmarshaler

type AdditionalFields interface {
	AdditionalField(key string, _ *jlexer.Lexer)
}
