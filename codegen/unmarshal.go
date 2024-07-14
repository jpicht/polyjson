package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
	"github.com/jpicht/polyjson/parser"
)

type UnmarshalFuncGen struct{}

func (UnmarshalFuncGen) GeneratePolyStruct(ctx *Context, p *generator.PolyStruct) error {
	ctx.Imports["github.com/mailru/easyjson/jlexer"] = "jlexer"
	ctx.Imports["github.com/jpicht/polyjson"] = "polyjson"

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
		fmt.Fprintf(ctx, "			ps.%s = new(%s)\n", impl.Struct.Name, impl.Struct.Name)
		if impl.Struct.Interfaces[parser.EasyJSONUnmarshaler] {
			fmt.Fprintf(ctx, "			ps.%s.UnmarshalEasyJSON(in)\n", impl.Struct.Name)
		} else if impl.Struct.Interfaces[parser.EncodingJSONUnmarshaler] {
			fmt.Fprintf(ctx, "			in.AddError(ps.%s.Unmarshal(in.Raw()))\n", impl.Struct.Name)
		} else {
			fmt.Fprintf(ctx, "			in.AddError(json.Unmarshal(in.Raw(), ps.%s))\n", impl.Struct.Name)
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
		fmt.Fprintf(ctx, "				Reason: \"unknown field\",\n")
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

	/*for _, impl := range p.Impls {
		fmt.Fprintf(w, "\tif ps.%s.Accept(v.Visit%s) {\n", impl.Struct.Name, impl.Struct.Name)
		fmt.Fprintf(w, "\t\treturn true\n")
		fmt.Fprintf(w, "\t}\n")
	}
	fmt.Fprintf(w, "\treturn false\n")*/
	fmt.Fprintf(ctx, "}\n\n")

	fmt.Fprintf(ctx, "func (ps *%s) UnmarshalJSON(data []byte) error {\n", p.Name)
	fmt.Fprintf(ctx, "	l := &jlexer.Lexer{Data: data}\n")
	fmt.Fprintf(ctx, "	ps.UnmarshalEasyJSON(l)\n")
	fmt.Fprintf(ctx, "	return l.Error()\n")
	fmt.Fprintf(ctx, "}\n\n")

	return nil
}
