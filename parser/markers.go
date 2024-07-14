package parser

import (
	_ "encoding/json"

	_ "github.com/mailru/easyjson"

	_ "github.com/jpicht/polyjson/markers"
	"github.com/jpicht/polyjson/parser/att"
)

var (
	MarkerImplements        = att.Mark("github.com/jpicht/polyjson/markers", "ImplementsInterface")
	MarkerCommon            = att.Mark("github.com/jpicht/polyjson/markers", "CommonInterface")
	MarkerInterface         = att.Mark("github.com/jpicht/polyjson/markers", "InterfaceInterface")
	AdditionalFields        = att.Mark("github.com/jpicht/polyjson/markers", "AdditionalFields")
	JWriterWritable         = att.Mark("github.com/launchdarkly/go-jsonstream/v3", "Writable")
	EasyJSONUnmarshaler     = att.Mark("github.com/mailru/easyjson", "Unmarshaler")
	EncodingJSONUnmarshaler = att.Mark("encoding/json", "Unmarshaler")
	EncodingJSONMarshaler   = att.Mark("encoding/json", "Marshaler")
)
