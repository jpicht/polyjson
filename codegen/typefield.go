package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type TypeFieldGen struct{}

func (TypeFieldGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	//fmt.Fprintf(w, "type %sType string\n\n", p.Name)
	if p.TypeID == nil {
		return nil
	}

	fmt.Fprintf(w, "const (\n")
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "\t%s%s = %s(%q)\n", p.TypeID.Name, impl.Struct.Name, p.TypeID.Name, TypeID(impl.Struct.Name, impl.Value.Tag))
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)
	return nil
}
