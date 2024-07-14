package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type AcceptFuncGen struct{}

func (AcceptFuncGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	fmt.Fprintf(w, "func (ps *%s) Accept(v %sVisitor) bool {\n", p.Name, p.Name)
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "\tif ps.%s != nil {\n", impl.Struct.Name)
		fmt.Fprintf(w, "\t\tv.Visit%s(*ps.%s)\n", impl.Struct.Name, impl.Struct.Name)
		fmt.Fprintf(w, "\t\treturn true\n")
		fmt.Fprintf(w, "\t}\n")
	}
	fmt.Fprintf(w, "\treturn false\n")
	fmt.Fprintf(w, "}\n\n")

	fmt.Fprintf(w, "func (pss %sSlice) Accept(v %sVisitor) bool {\n", p.Name, p.Name)
	fmt.Fprintf(w, "\tfor _, e := range pss {\n")
	fmt.Fprintf(w, "\t\tif !e.Accept(v) {\n")
	fmt.Fprintf(w, "\t\t\treturn false\n")
	fmt.Fprintf(w, "\t\t}\n")
	fmt.Fprintf(w, "\t}\n")
	fmt.Fprintf(w, "\treturn true\n")
	fmt.Fprintf(w, "}\n\n")
	return nil
}
