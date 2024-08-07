// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package testdata

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA96ca39cDecodeGithubComJpichtPolyjsonTestdata(in *jlexer.Lexer, out *CommonCD) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(true)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Timestamp":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA96ca39cEncodeGithubComJpichtPolyjsonTestdata(out *jwriter.Writer, in CommonCD) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Timestamp\":"
		out.RawString(prefix[1:])
		out.Raw((in.Timestamp).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CommonCD) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA96ca39cEncodeGithubComJpichtPolyjsonTestdata(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CommonCD) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA96ca39cEncodeGithubComJpichtPolyjsonTestdata(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CommonCD) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA96ca39cDecodeGithubComJpichtPolyjsonTestdata(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CommonCD) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA96ca39cDecodeGithubComJpichtPolyjsonTestdata(l, v)
}
