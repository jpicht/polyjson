package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type SliceTypeGen struct{}

func (SliceTypeGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	fmt.Fprintf(w, "type %sSlice []%s\n", p.Name, p.Name)
	return nil
}
