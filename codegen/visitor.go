package codegen

import (
	"fmt"

	"github.com/jpicht/polyjson/generator"
)

type VisitorInterfaceGen struct{}

func (VisitorInterfaceGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	fmt.Fprintf(w, "type %sVisitor interface {\n", p.Name)
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "\tVisit%s(%s)\n", impl.Struct.Name, impl.Struct.Name)
	}
	fmt.Fprintf(w, "}\n\n")
	return nil
}

type DefaultVisitorGen struct{}

func (DefaultVisitorGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	fmt.Fprintf(w, "type %sDefaultVisitor struct {}\n\n", p.Name)
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "func (%sDefaultVisitor) Visit%s(%s) {}\n", p.Name, impl.Struct.Name, impl.Struct.Name)
	}
	fmt.Fprintln(w)
	return nil
}

type VisitorFuncGen struct{}

func (VisitorFuncGen) GeneratePolyStruct(w *Context, p *generator.PolyStruct) error {
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "type %sVisitorFunc func(%s)\n", impl.Struct.Name, impl.Struct.Name)
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "type %sFuncVisitor struct {\n", p.Name)
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "\t%sVisitorFunc %sVisitorFunc\n", impl.Struct.Name, impl.Struct.Name)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	for _, impl := range p.Impls {
		fmt.Fprintf(w, "func (fv %sFuncVisitor) Visit%s(value %s) {\n", p.Name, impl.Struct.Name, impl.Struct.Name)
		fmt.Fprintf(w, "\tif fv.%sVisitorFunc != nil {\n", impl.Struct.Name)
		fmt.Fprintf(w, "\t\tfv.%sVisitorFunc(value)\n", impl.Struct.Name)
		fmt.Fprintln(w, "\t}")
		fmt.Fprintln(w, "}")
	}
	return nil
}
