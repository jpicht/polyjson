// Code generated by polyjson. DO NOT EDIT.
package testdata

import (
	"encoding/json"

	"github.com/jpicht/polyjson"
	"github.com/launchdarkly/go-jsonstream/v3/jwriter"
	"github.com/mailru/easyjson/jlexer"
)

type ABVisitor interface {
	VisitA(A)
	VisitB(B)
}

type AB struct {
	// common data
	CommonAB

	// implementations
	A polyjson.OneOf[A] `json:"a,omitempty"`
	B polyjson.OneOf[B] `json:"b,omitempty"`
}

type ABSlice []AB

func (ps *AB) Accept(v ABVisitor) bool {
	if ps.A.Accept(v.VisitA) {
		return true
	}
	if ps.B.Accept(v.VisitB) {
		return true
	}
	return false
}

func (pss ABSlice) Accept(v ABVisitor) bool {
	for _, e := range pss {
		if !e.Accept(v) {
			return false
		}
	}
	return true
}

func (ps *AB) WriteToJSONWriter(w *jwriter.Writer) {
	o := w.Object()
	defer o.End()

	// common fields from CommonAB
	var (
		raw []byte
		err error
	)
	raw, err = json.Marshal(ps.CommonAB.Timestamp)
	o.Maybe("timestamp", len(raw) > 0).Raw(raw)
	w.AddError(err)

	// implementations
	if ps.A.Valid() {
		ps.A.WriteToJSONWriter(o.Name("a"))
	} else if ps.B.Valid() {
		ps.B.WriteToJSONWriter(o.Name("b"))
	} else {
		w.AddError(polyjson.ErrNoValue)
	}
	// FIXME: additionals are not implemented yet
}

func (ps *AB) MarshalJSON() ([]byte, error) {
	w := jwriter.NewWriter()
	ps.WriteToJSONWriter(&w)
	return w.Bytes(), w.Error()
}

func (ps *AB) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}

	haveValue := false

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
		// implementations
		case "a":
			if haveValue {
				in.AddError(polyjson.ErrMultipleValues)
			}
			haveValue = true
			ps.A.UnmarshalEasyJSON(in)

		case "b":
			if haveValue {
				in.AddError(polyjson.ErrMultipleValues)
			}
			haveValue = true
			ps.B.UnmarshalEasyJSON(in)

		// common fields from CommonAB
		case "timestamp":
			in.AddError(ps.CommonAB.Timestamp.UnmarshalJSON(in.Raw()))

		// unknown fields are disallowed
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

	if !haveValue {
		in.AddError(polyjson.ErrNoValue)
	}
}

func (ps *AB) UnmarshalJSON(data []byte) error {
	l := &jlexer.Lexer{Data: data}
	ps.UnmarshalEasyJSON(l)
	return l.Error()
}

type ABDefaultVisitor struct{}

func (ABDefaultVisitor) VisitA(A) {}
func (ABDefaultVisitor) VisitB(B) {}

type (
	AVisitorFunc func(A)
	BVisitorFunc func(B)
)

type ABFuncVisitor struct {
	AVisitorFunc AVisitorFunc
	BVisitorFunc BVisitorFunc
}

func (fv ABFuncVisitor) VisitA(value A) {
	if fv.AVisitorFunc != nil {
		fv.AVisitorFunc(value)
	}
}

func (fv ABFuncVisitor) VisitB(value B) {
	if fv.BVisitorFunc != nil {
		fv.BVisitorFunc(value)
	}
}

type CDVisitor interface {
	VisitC(C)
	VisitD(D)
}

type CD struct {
	// common data
	CommonCD

	// implementations
	C polyjson.OneOf[C] `json:"c,omitempty"`
	D polyjson.OneOf[D] `json:"d,omitempty"`
}

type CDSlice []CD

func (ps *CD) Accept(v CDVisitor) bool {
	if ps.C.Accept(v.VisitC) {
		return true
	}
	if ps.D.Accept(v.VisitD) {
		return true
	}
	return false
}

func (pss CDSlice) Accept(v CDVisitor) bool {
	for _, e := range pss {
		if !e.Accept(v) {
			return false
		}
	}
	return true
}

func (ps *CD) WriteToJSONWriter(w *jwriter.Writer) {
	o := w.Object()
	defer o.End()

	// common fields from CommonCD
	var (
		raw []byte
		err error
	)
	raw, err = json.Marshal(ps.CommonCD.Timestamp)
	o.Maybe("timestamp", len(raw) > 0).Raw(raw)
	w.AddError(err)
	raw, err = json.Marshal(ps.CommonCD.SomeString)
	o.Maybe("some_string", len(raw) > 0).Raw(raw)
	w.AddError(err)
	raw, err = json.Marshal(ps.CommonCD.Additional)
	o.Maybe("additional", len(raw) > 0).Raw(raw)
	w.AddError(err)

	// implementations
	if ps.C.Valid() {
		ps.C.WriteToJSONWriter(o.Name("c"))
	} else if ps.D.Valid() {
		ps.D.WriteToJSONWriter(o.Name("d"))
	} else {
		w.AddError(polyjson.ErrNoValue)
	}
	// FIXME: additionals are not implemented yet
}

func (ps *CD) MarshalJSON() ([]byte, error) {
	w := jwriter.NewWriter()
	ps.WriteToJSONWriter(&w)
	return w.Bytes(), w.Error()
}

func (ps *CD) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}

	haveValue := false
	ps.CommonCD.Additional = make(polyjson.AdditionalFieldMap[json.RawMessage])

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
		// implementations
		case "c":
			if haveValue {
				in.AddError(polyjson.ErrMultipleValues)
			}
			haveValue = true
			ps.C.UnmarshalEasyJSON(in)

		case "d":
			if haveValue {
				in.AddError(polyjson.ErrMultipleValues)
			}
			haveValue = true
			ps.D.UnmarshalEasyJSON(in)

		// common fields from CommonCD
		case "timestamp":
			in.AddError(ps.CommonCD.Timestamp.UnmarshalJSON(in.Raw()))

		case "some_string":
			ps.CommonCD.SomeString = new(string)
			in.AddError(json.Unmarshal(in.Raw(), ps.CommonCD.SomeString))

		// unknown fields are allowed
		default:
			ps.CommonCD.Additional.AdditionalField(key, in)
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}

	if !haveValue {
		in.AddError(polyjson.ErrNoValue)
	}
}

func (ps *CD) UnmarshalJSON(data []byte) error {
	l := &jlexer.Lexer{Data: data}
	ps.UnmarshalEasyJSON(l)
	return l.Error()
}

type CDDefaultVisitor struct{}

func (CDDefaultVisitor) VisitC(C) {}
func (CDDefaultVisitor) VisitD(D) {}

type (
	CVisitorFunc func(C)
	DVisitorFunc func(D)
)

type CDFuncVisitor struct {
	CVisitorFunc CVisitorFunc
	DVisitorFunc DVisitorFunc
}

func (fv CDFuncVisitor) VisitC(value C) {
	if fv.CVisitorFunc != nil {
		fv.CVisitorFunc(value)
	}
}

func (fv CDFuncVisitor) VisitD(value D) {
	if fv.DVisitorFunc != nil {
		fv.DVisitorFunc(value)
	}
}

type EFVisitor interface {
	VisitE(E)
	VisitF(F)
}

type EF struct {
	// common data
	CommonEF

	// implementations
	E polyjson.OneOf[E] `json:"e,omitempty"`
	F polyjson.OneOf[F] `json:"f,omitempty"`
}

type EFSlice []EF

func (ps *EF) Accept(v EFVisitor) bool {
	if ps.E.Accept(v.VisitE) {
		return true
	}
	if ps.F.Accept(v.VisitF) {
		return true
	}
	return false
}

func (pss EFSlice) Accept(v EFVisitor) bool {
	for _, e := range pss {
		if !e.Accept(v) {
			return false
		}
	}
	return true
}

func (ps *EF) WriteToJSONWriter(w *jwriter.Writer) {
	o := w.Object()
	defer o.End()

	// implementations
	if ps.E.Valid() {
		ps.E.WriteToJSONWriter(o.Name("e"))
	} else if ps.F.Valid() {
		ps.F.WriteToJSONWriter(o.Name("f"))
	} else {
		w.AddError(polyjson.ErrNoValue)
	}
	// FIXME: additionals are not implemented yet
}

func (ps *EF) MarshalJSON() ([]byte, error) {
	w := jwriter.NewWriter()
	ps.WriteToJSONWriter(&w)
	return w.Bytes(), w.Error()
}

func (ps *EF) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}

	haveValue := false

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
		// implementations
		case "e":
			if haveValue {
				in.AddError(polyjson.ErrMultipleValues)
			}
			haveValue = true
			ps.E.UnmarshalEasyJSON(in)

		case "f":
			if haveValue {
				in.AddError(polyjson.ErrMultipleValues)
			}
			haveValue = true
			ps.F.UnmarshalEasyJSON(in)

		// unknown fields are allowed
		default:
			ps.CommonEF.AdditionalField(key, in)
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}

	if !haveValue {
		in.AddError(polyjson.ErrNoValue)
	}
}

func (ps *EF) UnmarshalJSON(data []byte) error {
	l := &jlexer.Lexer{Data: data}
	ps.UnmarshalEasyJSON(l)
	return l.Error()
}

type EFDefaultVisitor struct{}

func (EFDefaultVisitor) VisitE(E) {}
func (EFDefaultVisitor) VisitF(F) {}

type (
	EVisitorFunc func(E)
	FVisitorFunc func(F)
)

type EFFuncVisitor struct {
	EVisitorFunc EVisitorFunc
	FVisitorFunc FVisitorFunc
}

func (fv EFFuncVisitor) VisitE(value E) {
	if fv.EVisitorFunc != nil {
		fv.EVisitorFunc(value)
	}
}

func (fv EFFuncVisitor) VisitF(value F) {
	if fv.FVisitorFunc != nil {
		fv.FVisitorFunc(value)
	}
}
