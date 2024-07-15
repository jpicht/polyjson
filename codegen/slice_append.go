package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type SliceAppendGen struct{}

func (SliceAppendGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	var (
		commonParam = ""
		commonField = ""
	)

	if p.Common != nil {
		commonField = "\t\t" + p.Common.Name + ": c,\n"
		commonParam = "c " + p.Common.Name + ","
	}

	fmt.Fprintf(w, "func (s *%sSlice) Append(value %s) {\n", p.Name, p.Name)
	fmt.Fprintf(w, "\t*s = append(*s, value)\n")
	fmt.Fprintf(w, "}\n\n")

	// TODO: type ID field

	for _, impl := range p.Impls {
		fmt.Fprintf(w, "func (s *%sSlice) Append%s(%svalue %s) {\n", p.Name, impl.Struct.Name, commonParam, impl.Struct.Name)
		fmt.Fprintf(w, "\t*s = append(*s, %s{\n", p.Name)
		fmt.Fprint(w, commonField)
		fmt.Fprintf(w, "\t\t%s: &value,\n", impl.Struct.Name)
		fmt.Fprintf(w, "\t})\n")
		fmt.Fprintf(w, "}\n\n")
	}

	return nil
}
