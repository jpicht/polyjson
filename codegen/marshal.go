package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
	"github.com/jpicht/polyjson/parser"
)

type MarshalFuncGen struct{}

func (MarshalFuncGen) GeneratePolyStruct(ctx *Context, p *generator.PolyStruct) error {
	ctx.Imports["github.com/launchdarkly/go-jsonstream/v3/jwriter"] = "jwriter"
	ctx.Imports["github.com/jpicht/polyjson"] = "polyjson"

	fmt.Fprintf(ctx, "func (ps *%s) WriteToJSONWriter(w*jwriter.Writer) {\n", p.Name)
	fmt.Fprintln(ctx, "	o := w.Object()")
	fmt.Fprintln(ctx, "	defer o.End()")
	fmt.Fprintln(ctx, "")

	if p.Common != nil && len(p.Common.Fields) > 0 {
		fmt.Fprintf(ctx, "	// common fields from %s\n", p.Common.Name)
		fmt.Fprintln(ctx, "	var (")
		fmt.Fprintln(ctx, "		raw []byte")
		fmt.Fprintln(ctx, "		err error")
		fmt.Fprintln(ctx, "	)")
		for _, field := range p.Common.Fields {
			fmt.Fprintf(ctx, "	raw, err = json.Marshal(ps.%s.%s)\n", p.Common.Name, field.Name)
			fmt.Fprintf(ctx, "	o.Maybe(%q, len(raw) > 0).Raw(raw)\n", JSONName(field.Name, field.Tag))
			fmt.Fprintf(ctx, "	w.AddError(err)\n")
		}
	}
	fmt.Fprintln(ctx, "")

	fmt.Fprintf(ctx, "	// implementations\n")
	elseStr := ""
	for _, impl := range p.Impls {
		if impl.Struct.Interfaces[parser.JWriterWritable] {
			fmt.Fprintf(ctx, "	%sif value, ok := ps.%s.Get(); ok {\n", elseStr, impl.Struct.Name)
			fmt.Fprintf(ctx, "		value.WriteToJSONWriter(o.Name(%q))\n", JSONName(impl.Struct.Name, impl.Value.Tag))
		}
		fmt.Fprintf(ctx, "	%sif ps.%s.Valid() {\n", elseStr, impl.Struct.Name)
		fmt.Fprintf(ctx, "		ps.%s.WriteToJSONWriter(o.Name(%q))\n", impl.Struct.Name, JSONName(impl.Struct.Name, impl.Value.Tag))
		elseStr = "} else "
	}
	fmt.Fprintln(ctx, "	} else {")
	fmt.Fprintln(ctx, "		w.AddError(polyjson.ErrNoValue)")
	fmt.Fprintln(ctx, "	}")
	fmt.Fprintf(ctx, "		// FIXME: additionals are not implemented yet\n")
	fmt.Fprintf(ctx, "}\n\n")

	fmt.Fprintf(ctx, "func (ps *%s) MarshalJSON() ([]byte, error) {\n", p.Name)
	fmt.Fprintf(ctx, "	w := jwriter.NewWriter()\n")
	fmt.Fprintf(ctx, "	ps.WriteToJSONWriter(&w)\n")
	fmt.Fprintf(ctx, "	return w.Bytes(), w.Error()\n")
	fmt.Fprintf(ctx, "}\n\n")

	return nil
}
