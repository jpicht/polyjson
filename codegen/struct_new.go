package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type PolyStructNewGen struct{}

func (PolyStructNewGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	var (
		commonParam = ""
		commonField = ""
		commonCopy  = ""
	)

	fmt.Fprintf(w, "type %sBuilder struct {\n", p.Name)
	if p.Common != nil {
		fmt.Fprintf(w, "\t%s\n\n", p.Common.Name)
		commonField = p.Common.Name + ": c"
		commonParam = "c " + p.Common.Name
		commonCopy = "\t\t" + p.Common.Name + ": b." + p.Common.Name + ",\n"
	}
	fmt.Fprintf(w, "}\n\n")

	fmt.Fprintf(w, "func New%s(%s) %sBuilder {\n", p.Name, commonParam, p.Name)
	fmt.Fprintf(w, "\treturn %sBuilder{ %s }\n\n", p.Name, commonField)
	fmt.Fprintf(w, "}\n\n")

	// TODO: type ID field

	for _, impl := range p.Impls {
		fmt.Fprintf(w, "func (b %sBuilder) %s(value %s) %s {\n", p.Name, impl.Struct.Name, impl.Struct.Name, p.Name)
		fmt.Fprintf(w, "\treturn %s{\n", p.Name)
		fmt.Fprint(w, commonCopy)
		fmt.Fprintf(w, "\t\t%s: &value,\n", impl.Struct.Name)
		fmt.Fprintf(w, "\t}\n")
		fmt.Fprintf(w, "}\n\n")
	}

	return nil
}
