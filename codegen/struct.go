package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type PolyStructGen struct{}

func (PolyStructGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	fmt.Fprintf(w, "type %s struct {\n", p.Name)
	if p.Common != nil {
		fmt.Fprintf(w, "\t// common data\n")
		fmt.Fprintf(w, "\t%s\n\n", p.Common.Name)
	}
	fmt.Fprintf(w, "\t// implementations\n")
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "\t%s polyjson.OneOf[%s] `json:\"%s\"`\n", impl.Struct.Name, impl.Struct.Name, JSONTag(impl.Struct.Name, impl.Value.Tag))
	}
	fmt.Fprintf(w, "}\n\n")
	return nil
}
