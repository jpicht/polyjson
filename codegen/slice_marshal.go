package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type SliceMarshalGen struct{}

func (SliceMarshalGen) GeneratePolyStruct(ctx *Context, p *generator.PolyStruct) error {
	ctx.Imports["github.com/launchdarkly/go-jsonstream/v3/jwriter"] = "jwriter"

	fmt.Fprintf(ctx, "func (s *%sSlice) WriteToJSONWriter(w *jwriter.Writer) {\n", p.Name)
	fmt.Fprintf(ctx, "\ta := w.Array()\n")
	fmt.Fprintf(ctx, "\tdefer a.End()\n")
	fmt.Fprintf(ctx, "\tfor _, item := range *s {\n")
	fmt.Fprintf(ctx, "\t\titem.WriteToJSONWriter(w)\n")
	fmt.Fprintf(ctx, "\t}\n")
	fmt.Fprintf(ctx, "}\n\n")

	fmt.Fprintf(ctx, "func (s *%sSlice) MarshalJSON() ([]byte, error) {\n", p.Name)
	fmt.Fprintf(ctx, "\tw := jwriter.NewWriter()\n")
	fmt.Fprintf(ctx, "\ts.WriteToJSONWriter(&w)\n")
	fmt.Fprintf(ctx, "\treturn w.Bytes(), w.Error()\n")
	fmt.Fprintf(ctx, "}\n\n")
	return nil
}
