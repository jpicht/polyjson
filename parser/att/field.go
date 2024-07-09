package att

import (
	"go/types"
	"reflect"
)

type Field struct {
	Struct     *NamedStruct
	Name       string
	Type       Type
	Tag        reflect.StructTag
	Interfaces map[Marker]bool
}

func FieldIs[T Type](f *Field) bool {
	return TypeIs[T](f.Type)
}

func TypeIs[T Type](t Type) bool {
	if _, ok := t.(T); ok {
		return true
	}

	switch tt := t.(type) {
	case *types.Alias:
		return TypeIs[T](tt.Obj().Type())

	case *types.Named:
		return TypeIs[T](tt.Obj().Type())

	default:
		return false
	}
}

func (f *Field) IsPointer() bool {
	for l, t := types.Type(nil), f.Type; t != l; l, t = t, t.Underlying() {
		if _, ok := t.(*types.Pointer); ok {
			return true
		}
	}
	return false
}

func (f *Field) IsMap() bool {
	for l, t := types.Type(nil), f.Type; t != l; l, t = t, t.Underlying() {
		if _, ok := t.(*types.Map); ok {
			return true
		}
	}
	return false
}
