package parser

import (
	_ "encoding/json"

	_ "github.com/mailru/easyjson"

	"github.com/jpicht/polyjson/parser/att"
)

var (
	MarkerImplements        = att.Mark("github.com/jpicht/polyjson", "Implements", "ImplementsInterface")
	MarkerCommon            = att.Mark("github.com/jpicht/polyjson", "Common", "CommonInterface")
	MarkerInterface         = att.Mark("github.com/jpicht/polyjson", "Interface", "InterfaceInterface")
	AdditionalFields        = att.Mark("github.com/jpicht/polyjson", "", "AdditionalFields")
	JWriterWritable         = att.Mark("github.com/launchdarkly/go-jsonstream/v3", "", "Writable")
	EasyJSONUnmarshaler     = att.Mark("github.com/mailru/easyjson", "", "Unmarshaler")
	EncodingJSONUnmarshaler = att.Mark("encoding/json", "", "Unmarshaler")
	EncodingJSONMarshaler   = att.Mark("encoding/json", "", "Marshaler")
)
