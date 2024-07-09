package codegen

import (
	"github.com/jpicht/polyjson/generator"
)

type CodeGen interface {
	GeneratePolyStruct(*Context, *generator.PolyStruct) error
}
