package parser

import (
	_ "encoding/json"

	_ "github.com/mailru/easyjson"

	"github.com/jpicht/polyjson/parser/att"
)

var (
	MarkerIs                = att.Mark("github.com/jpicht/polyjson", "Is", "IsInterface")
	MarkerCommon            = att.Mark("github.com/jpicht/polyjson", "Common", "CommonInterface")
	AdditionalFields        = att.Mark("github.com/jpicht/polyjson", "", "AdditionalFields")
	JWriterWritable         = att.Mark("github.com/launchdarkly/go-jsonstream/v3", "", "Writable")
	EasyJSONUnmarshaler     = att.Mark("github.com/mailru/easyjson", "", "Unmarshaler")
	EncodingJSONUnmarshaler = att.Mark("encoding/json", "", "Unmarshaler")
)
