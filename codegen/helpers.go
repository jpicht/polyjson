package codegen

import (
	"reflect"
	"strings"

	"github.com/stoewer/go-strcase"
)

func JSONName(name string, tag reflect.StructTag) string {
	n, _, _ := strings.Cut(JSONTag(name, tag), ",")
	if n != "" {
		return n
	}
	return strcase.SnakeCase(name)
}

func JSONTag(name string, tag reflect.StructTag) string {
	if v, ok := tag.Lookup("json"); ok {
		return v
	}
	if v, ok := tag.Lookup("polyjson"); ok {
		return v
	}
	return strcase.SnakeCase(name) + ",omitempty"
}

func TypeID(name string, tag reflect.StructTag) string {
	if v, ok := tag.Lookup("polytypeid"); ok {
		return v
	}
	return strcase.SnakeCase(name)
}
