package markers

import (
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type Interface interface{ isInterface() }
type Implements interface{ isImplementation() }
type Common interface{ isCommon() }

type EasyJSONUnmarshaller = easyjson.Unmarshaler

type AdditionalFields interface {
	AdditionalField(key string, _ *jlexer.Lexer)
}
