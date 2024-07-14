package att

import (
	"go/types"
	"reflect"
)

type Type = types.Type

type Named[Underlying Type] struct {
	*types.Named
	TypedUnderlying Type
}

type NamedStruct struct {
	File  string
	Name  string
	Named *types.Named
	*types.Struct
	Fields     []*Field
	Markers    map[Marker]MarkerValue
	Interfaces map[Marker]bool
}

type MarkerValue struct {
	Target string
	Field  *Field
	Tag    reflect.StructTag
}

func (n *NamedStruct) Is(m Marker) bool {
	_, ok := n.Markers[m]
	return ok
}

func (mv MarkerValue) String() string {
	return "[-> " + mv.Target + "]"
}

func (n *NamedStruct) FindField(tm TypeMarker) *Field {
	if n == nil {
		return nil
	}
	for _, f := range n.Fields {
		if f.Type == tm.Type {
			return f
		}
	}
	return nil
}
