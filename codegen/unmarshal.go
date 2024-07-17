package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/jpicht/polyjson/generator"
	"github.com/jpicht/polyjson/parser"
	"github.com/jpicht/polyjson/parser/att"
)

type UnmarshalFuncGen struct{}

func (u UnmarshalFuncGen) GeneratePolyStruct(ctx *Context, p *generator.PolyStruct) error {
	if p.TypeID != nil {
		return u.withTypeField(ctx, p)
	}
	return u.withoutTypeField(ctx, p)
}

func (u UnmarshalFuncGen) withTypeField(ctx *Context, p *generator.PolyStruct) error {
	ctx.Imports["github.com/mailru/easyjson/jlexer"] = "jlexer"
	ctx.Imports["github.com/jpicht/polyjson"] = "polyjson"
	ctx.Imports["encoding/json"] = "json"
	ctx.Imports["fmt"] = "fmt"

	if p.Additional != nil {
		return fmt.Errorf("%s combines additional fields and type ID", p.Name)
	}

	var (
		typeIDField   *att.Field
		typeFieldName = "ps.Type"
	)

	if p.Common != nil {
		typeIDField = p.Common.FindField(*p.TypeID)
		typeFieldName = fmt.Sprintf("ps.%s.%s", p.Common.Name, typeIDField.Name)
	}

	fmt.Fprintf(ctx, "func (ps *%s) UnmarshalEasyJSON(in *jlexer.Lexer) {\n", p.Name)
	fmt.Fprintf(ctx, "	isTopLevel := in.IsStart()\n")
	fmt.Fprintf(ctx, "	if in.IsNull() {\n")
	fmt.Fprintf(ctx, "		if isTopLevel {\n")
	fmt.Fprintf(ctx, "			in.Consumed()\n")
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		in.Skip()\n")
	fmt.Fprintf(ctx, "		return\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintln(ctx)

	fmt.Fprintf(ctx, "	cache := make(map[string]json.RawMessage)\n")
	fmt.Fprintln(ctx)

	fmt.Fprintf(ctx, "	in.Delim('{')\n")
	fmt.Fprintf(ctx, "	for !in.IsDelim('}') {\n")
	fmt.Fprintf(ctx, "		key := in.UnsafeFieldName(true)\n")
	fmt.Fprintf(ctx, "		in.WantColon()\n")
	fmt.Fprintf(ctx, "		if in.IsNull() {\n")
	fmt.Fprintf(ctx, "			in.Skip()\n")
	fmt.Fprintf(ctx, "			in.WantComma()\n")
	fmt.Fprintf(ctx, "			continue\n")
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		switch key {\n")

	if typeIDField == nil {
		fmt.Fprintf(ctx, "		// type ID\n")
		fmt.Fprintf(ctx, "		case \"type\":\n")
		fmt.Fprintf(ctx, "			ps.Type = %s(in.String())\n", p.TypeID.Name)
		fmt.Fprintln(ctx)
	}

	if err := u.commonFields(ctx, p); err != nil {
		return err
	}

	fmt.Fprintf(ctx, "		// cache all fields\n")
	fmt.Fprintf(ctx, "		default:\n")
	fmt.Fprintf(ctx, "			cache[key] = in.Raw()\n")
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		in.WantComma()\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintf(ctx, "	in.Delim('}')\n")
	fmt.Fprintf(ctx, "	if isTopLevel {\n")
	fmt.Fprintf(ctx, "		in.Consumed()\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintln(ctx)

	fmt.Fprintf(ctx, "	raw, err := json.Marshal(cache)\n")
	fmt.Fprintf(ctx, "	in.AddError(err)\n")
	fmt.Fprintf(ctx, "	cache = nil\n")
	fmt.Fprintln(ctx)

	fmt.Fprintf(ctx, "	switch %s {\n", typeFieldName)
	for _, impl := range p.Impls {
		fmt.Fprintf(ctx, "	case %s%s:\n", p.TypeID.Name, impl.Struct.Name)

		fmt.Fprintf(ctx, "		ps.%s = &%s{\n", impl.Struct.Name, impl.Struct.Name)
		fmt.Fprintf(ctx, "			Implements: polyjson.Implements[%s]{ Parent: ps },\n", p.Name)
		fmt.Fprintf(ctx, "		}\n")

		if impl.Struct.Interfaces[parser.EasyJSONUnmarshaler] {
			fmt.Fprintf(ctx, "		ps.%s.UnmarshalEasyJSON(&jlexer.Lexer{ Data: raw })\n", impl.Struct.Name)
		} else if impl.Struct.Interfaces[parser.EncodingJSONUnmarshaler] {
			fmt.Fprintf(ctx, "		in.AddError(ps.%s.Unmarshal(raw))\n", impl.Struct.Name)
		} else {
			fmt.Fprintf(ctx, "		in.AddError(json.Unmarshal(raw, ps.%s))\n", impl.Struct.Name)
		}
		fmt.Fprintln(ctx)
	}

	fmt.Fprintf(ctx, "	case \"\":\n")
	fmt.Fprintf(ctx, "		in.AddError(polyjson.ErrMissingTypeID)\n")
	fmt.Fprintf(ctx, "		return\n")
	fmt.Fprintf(ctx, "	default:\n")
	fmt.Fprintf(ctx, "		in.AddError(fmt.Errorf(\"%%w: %%s\", polyjson.ErrInvalidTypeID, %s))\n", typeFieldName)
	fmt.Fprintf(ctx, "		return\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintf(ctx, "}\n\n")

	fmt.Fprintf(ctx, "func (ps *%s) UnmarshalJSON(data []byte) error {\n", p.Name)
	fmt.Fprintf(ctx, "	l := &jlexer.Lexer{Data: data}\n")
	fmt.Fprintf(ctx, "	ps.UnmarshalEasyJSON(l)\n")
	fmt.Fprintf(ctx, "	return l.Error()\n")
	fmt.Fprintf(ctx, "}\n\n")

	return nil
}

func (UnmarshalFuncGen) commonFields(ctx *Context, p *generator.PolyStruct) error {
	if p.Common == nil || len(p.Common.Fields) < 1 {
		return nil
	}

	fmt.Fprintf(ctx, "		// common fields from %s\n", p.Common.Name)
	for _, field := range p.Common.Fields {
		addrOp := "&"
		render := func() {
			fmt.Fprintf(ctx, "		case %q:\n", JSONName(field.Name, field.Tag))
			if field.IsPointer() {
				fmt.Fprintf(ctx, "			ps.%s.%s = %s\n", p.Common.Name, field.Name, ctx.New(field.Type))
				addrOp = ""
			}
		}
		switch {
		case field.Interfaces[parser.AdditionalFields]:
			p.Additional = field
		case field.Interfaces[parser.EasyJSONUnmarshaler]:
			render()
			fmt.Fprintf(ctx, "			ps.%s.%s.UnmarshalEasyJSON(in)\n", p.Common.Name, field.Name)
		case field.Interfaces[parser.EncodingJSONUnmarshaler]:
			render()
			fmt.Fprintf(ctx, "			in.AddError(ps.%s.%s.UnmarshalJSON(in.Raw()))\n", p.Common.Name, field.Name)
		default:
			render()
			ctx.Imports["encoding/json"] = "json"
			fmt.Fprintf(ctx, "			in.AddError(json.Unmarshal(in.Raw(), %sps.%s.%s))\n", addrOp, p.Common.Name, field.Name)
		}
		fmt.Fprintln(ctx)
	}

	return nil
}

func (UnmarshalFuncGen) withoutTypeField(ctx *Context, p *generator.PolyStruct) error {
	ctx.Imports["github.com/mailru/easyjson/jlexer"] = "jlexer"
	ctx.Imports["github.com/jpicht/polyjson"] = "polyjson"
	ctx.Imports["fmt"] = "fmt"

	fmt.Fprintf(ctx, "func (ps *%s) UnmarshalEasyJSON(in *jlexer.Lexer) {\n", p.Name)
	fmt.Fprintf(ctx, "	isTopLevel := in.IsStart()\n")
	fmt.Fprintf(ctx, "	if in.IsNull() {\n")
	fmt.Fprintf(ctx, "		if isTopLevel {\n")
	fmt.Fprintf(ctx, "			in.Consumed()\n")
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		in.Skip()\n")
	fmt.Fprintf(ctx, "		return\n")
	fmt.Fprintf(ctx, "	}\n")

	fmt.Fprintln(ctx)
	fmt.Fprintf(ctx, "	haveValue := false\n")
	if p.Additional != nil && (p.Additional.IsPointer() || p.Additional.IsMap()) {
		fmt.Fprintf(ctx, "	ps.%s.%s = %s\n", p.Common.Name, p.Additional.Name, ctx.New(p.Additional.Type))
	}
	fmt.Fprintln(ctx)

	needDecodeBuffer := false

	for _, impl := range p.Impls {
		if impl.Struct.Interfaces[parser.EasyJSONUnmarshaler] {
			continue
		}
		if impl.Struct.Interfaces[parser.EncodingJSONUnmarshaler] {
			continue
		}
		needDecodeBuffer = true
	}

	if needDecodeBuffer {
		ctx.Imports["bytes"] = "bytes"
		fmt.Fprintf(ctx, "	buf := &bytes.Buffer{}\n")
		fmt.Fprintf(ctx, "	dec := json.NewDecoder(buf)\n")
		fmt.Fprintf(ctx, "	dec.DisallowUnknownFields()\n")
	}

	fmt.Fprintf(ctx, "	wrapErr := func(err error, field string) error {\n")
	fmt.Fprintf(ctx, "		if err != nil {\n")
	fmt.Fprintf(ctx, "			err = fmt.Errorf(\"error parsing %s.%%s: %%w\", field, err)\n", p.Name)
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		return err\n")
	fmt.Fprintf(ctx, "	}\n")

	fmt.Fprintf(ctx, "	in.Delim('{')\n")
	fmt.Fprintf(ctx, "	for !in.IsDelim('}') {\n")
	fmt.Fprintf(ctx, "		key := in.UnsafeFieldName(true)\n")
	fmt.Fprintf(ctx, "		in.WantColon()\n")
	fmt.Fprintf(ctx, "		if in.IsNull() {\n")
	fmt.Fprintf(ctx, "			in.Skip()\n")
	fmt.Fprintf(ctx, "			in.WantComma()\n")
	fmt.Fprintf(ctx, "			continue\n")
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		switch key {\n")

	fmt.Fprintf(ctx, "		// implementations\n")
	for _, impl := range p.Impls {
		fmt.Fprintf(ctx, "		case %q:\n", JSONName(impl.Struct.Name, impl.Value.Tag))
		fmt.Fprintf(ctx, "			if haveValue {\n")
		fmt.Fprintf(ctx, "				in.AddError(polyjson.ErrMultipleValues)\n")
		fmt.Fprintf(ctx, "			}\n")
		fmt.Fprintf(ctx, "			haveValue = true\n")

		fmt.Fprintf(ctx, "			ps.%s = &%s{\n", impl.Struct.Name, impl.Struct.Name)
		fmt.Fprintf(ctx, "				Implements: polyjson.Implements[%s]{ Parent: ps },\n", p.Name)
		fmt.Fprintf(ctx, "			}\n")

		if impl.Struct.Interfaces[parser.EasyJSONUnmarshaler] {
			fmt.Fprintf(ctx, "			ps.%s.UnmarshalEasyJSON(in)\n", impl.Struct.Name)
		} else if impl.Struct.Interfaces[parser.EncodingJSONUnmarshaler] {
			fmt.Fprintf(ctx, "			in.AddError(wrapErr(ps.%s.Unmarshal(in.Raw()), %q))\n", impl.Struct.Name, impl.Struct.Name)
		} else {
			fmt.Fprintf(ctx, "			buf.Write(in.Raw())\n")
			fmt.Fprintf(ctx, "			in.AddError(wrapErr(dec.Decode(ps.%s), %q))\n", impl.Struct.Name, impl.Struct.Name)
			fmt.Fprintf(ctx, "			buf.Reset()\n")
		}

		fmt.Fprintln(ctx)
	}

	if p.Common != nil && len(p.Common.Fields) > 0 {
		fmt.Fprintf(ctx, "		// common fields from %s\n", p.Common.Name)
		for _, field := range p.Common.Fields {
			addrOp := "&"
			render := func() {
				fmt.Fprintf(ctx, "		case %q:\n", JSONName(field.Name, field.Tag))
				if field.IsPointer() {
					fmt.Fprintf(ctx, "			ps.%s.%s = %s\n", p.Common.Name, field.Name, ctx.New(field.Type))
					addrOp = ""
				}
			}
			switch {
			case field.Interfaces[parser.AdditionalFields]:
				p.Additional = field
			case field.Interfaces[parser.EasyJSONUnmarshaler]:
				render()
				fmt.Fprintf(ctx, "			ps.%s.%s.UnmarshalEasyJSON(in)\n", p.Common.Name, field.Name)
			case field.Interfaces[parser.EncodingJSONUnmarshaler]:
				render()
				fmt.Fprintf(ctx, "			in.AddError(ps.%s.%s.UnmarshalJSON(in.Raw()))\n", p.Common.Name, field.Name)
			default:
				render()
				ctx.Imports["encoding/json"] = "json"
				fmt.Fprintf(ctx, "			in.AddError(json.Unmarshal(in.Raw(), %sps.%s.%s))\n", addrOp, p.Common.Name, field.Name)
			}
			fmt.Fprintln(ctx)
		}
	}
	if p.Additional != nil {
		fmt.Fprintf(ctx, "		// unknown fields are allowed\n")
		fmt.Fprintf(ctx, "		default:\n")
		fmt.Fprintf(ctx, "			ps.%s.%s.AdditionalField(key, in)\n", p.Common.Name, p.Additional.Name)
	} else if p.Common != nil && p.Common.Interfaces[parser.AdditionalFields] {
		fmt.Fprintf(ctx, "		// unknown fields are allowed\n")
		fmt.Fprintf(ctx, "		default:\n")
		fmt.Fprintf(ctx, "			ps.%s.AdditionalField(key, in)\n", p.Common.Name)
	} else {
		fmt.Fprintf(ctx, "		// unknown fields are disallowed\n")
		fmt.Fprintf(ctx, "		default:\n")
		fmt.Fprintf(ctx, "			in.AddError(&jlexer.LexerError{\n")
		fmt.Fprintf(ctx, "				Offset: in.GetPos(),\n")
		fmt.Fprintf(ctx, "				Reason: \"unknown field in %s.%s\",\n", filepath.Base(p.Package.Name), p.Name)
		fmt.Fprintf(ctx, "				Data:   key,\n")
		fmt.Fprintf(ctx, "			})\n")
	}
	fmt.Fprintf(ctx, "		}\n")
	fmt.Fprintf(ctx, "		in.WantComma()\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintf(ctx, "	in.Delim('}')\n")
	fmt.Fprintf(ctx, "	if isTopLevel {\n")
	fmt.Fprintf(ctx, "		in.Consumed()\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintln(ctx)
	fmt.Fprintf(ctx, "	if !haveValue {\n")
	fmt.Fprintf(ctx, "		in.AddError(polyjson.ErrNoValue)\n")
	fmt.Fprintf(ctx, "	}\n")
	fmt.Fprintf(ctx, "}\n\n")

	fmt.Fprintf(ctx, "func (ps *%s) UnmarshalJSON(data []byte) error {\n", p.Name)
	fmt.Fprintf(ctx, "	l := &jlexer.Lexer{Data: data}\n")
	fmt.Fprintf(ctx, "	ps.UnmarshalEasyJSON(l)\n")
	fmt.Fprintf(ctx, "	return l.Error()\n")
	fmt.Fprintf(ctx, "}\n\n")

	return nil
}
