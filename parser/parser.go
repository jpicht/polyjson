package parser

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"reflect"
	"strings"

	"github.com/jpicht/polyjson/parser/att"
	"golang.org/x/tools/go/packages"
)

type Result struct {
	Config *Config
	StdLib struct {
		FileSet  *token.FileSet
		Packages []*packages.Package
	}
	Additional map[string][]byte
	Directory  string
	Packages   []*att.Package
	Helpers    struct {
		InterfaceTests map[att.Marker]func(types.Type) bool
		MarkerTests    map[att.Marker]func(types.Type) bool
	}
}

func (config *Config) Parse(directory string) (*Result, error) {
	config.Packages.Dir = directory

	r := &Result{
		Config:     config,
		Directory:  directory,
		Additional: make(map[string][]byte),
	}

	r.StdLib.FileSet = token.NewFileSet()

	r.Helpers.MarkerTests = make(map[att.Marker]func(types.Type) bool, len(config.Markers))
	r.Helpers.InterfaceTests = make(map[att.Marker]func(types.Type) bool, len(config.Interfaces))

	if err := r.load(); err != nil {
		return nil, err
	}

	for _, marker := range config.Markers {
		if _, ok := r.Helpers.MarkerTests[marker]; ok {
			continue
		}
		for _, inPkg := range r.StdLib.Packages {
			pkg, ok := inPkg.Imports[marker.Package]
			if !ok {
				continue
			}
			markerInterface := pkg.Types.Scope().Lookup(marker.InterfaceName)
			if markerInterface != nil {
				r.Helpers.MarkerTests[marker] = makeMarkerCheck(markerInterface.Type())
			}
		}
	}

	for _, iface := range config.Interfaces {
		if _, ok := r.Helpers.InterfaceTests[iface]; ok {
			continue
		}
		for _, inPkg := range r.StdLib.Packages {
			pkg, ok := inPkg.Imports[iface.Package]
			if !ok {
				continue
			}
			markerInterface := pkg.Types.Scope().Lookup(iface.InterfaceName)
			if markerInterface != nil {
				r.Helpers.InterfaceTests[iface] = makeMarkerCheck(markerInterface.Type())
			}
		}
	}

	for _, inPkg := range r.StdLib.Packages {
		outPkg, err := r.parsePackage(inPkg)
		if err != nil {
			return nil, fmt.Errorf("error in package %q: %w", inPkg.ID, err)
		}
		r.Packages = append(r.Packages, outPkg)
	}

	return r, nil
}

func makeMarkerCheck(iface types.Type) func(types.Type) bool {
	return func(t types.Type) bool {
		return types.AssignableTo(t, iface) ||
			types.AssignableTo(types.NewPointer(t), iface)
	}
}

func (r *Result) ignored(pos token.Pos) bool {
	file := r.StdLib.FileSet.File(pos)
	return strings.HasSuffix(file.Name(), "_polyjson.go")
}

func (r *Result) parsePackage(inPkg *packages.Package) (*att.Package, error) {
	log.Printf("parsing package %q", inPkg.ID)

	qual := types.RelativeTo(inPkg.Types)
	scope := inPkg.Types.Scope()

	outPkg := att.Package{
		Package: inPkg,
	}

	for _, name := range scope.Names() {
		obj := scope.Lookup(name)

		named, structT, ok := getStruct(obj.Type())
		if !ok {
			continue
		}

		outStruct := &att.NamedStruct{
			File:       r.StdLib.FileSet.File(obj.Pos()).Name(),
			Name:       types.TypeString(named, qual),
			Named:      named,
			Struct:     structT,
			Markers:    make(map[att.Marker]att.MarkerValue),
			Interfaces: make(map[att.Marker]bool),
		}

		for i := 0; i < structT.NumFields(); i++ {
			field := structT.Field(i)

			f := &att.Field{
				Struct:     outStruct,
				Tag:        reflect.StructTag(structT.Tag(i)),
				Name:       field.Name(),
				Type:       field.Type(),
				Interfaces: make(map[att.Marker]bool),
			}

			isMarker := false
			for marker, testFn := range r.Helpers.MarkerTests {
				if testFn(field.Type()) {
					t, ok := getFirstTypeArg[*types.Named](field.Type())
					if !ok {
						//	log.Printf("invalid type argument len %d, expected 1", args.Len())
						log.Printf("skipping field %#v, no type arg", field)
						continue
					}
					target := types.TypeString(t, qual)
					outStruct.Markers[marker] = att.MarkerValue{
						Target: target,
						Tag:    reflect.StructTag(structT.Tag(i)),
						Field:  f,
					}
					isMarker = true
				}
			}

			if isMarker {
				continue
			}

			for iface, testFn := range r.Helpers.InterfaceTests {
				if !testFn(f.Type) {
					continue
				}
				f.Interfaces[iface] = true
			}

			outStruct.Fields = append(outStruct.Fields, f)
		}

		for iface, testFn := range r.Helpers.InterfaceTests {
			if !testFn(named) {
				continue
			}
			outStruct.Interfaces[iface] = true
		}

		outPkg.NamedStructs = append(outPkg.NamedStructs, outStruct)
	}

	return &outPkg, nil
}

func getStruct(t types.Type) (named *types.Named, structT *types.Struct, ok bool) {
	named, ok = t.(*types.Named)
	if !ok {
		return
	}

	structT, ok = named.Underlying().(*types.Struct)
	return
}

func getTypeArgs(t types.Type) *types.TypeList {
	switch typed := t.(type) {
	case *types.Named:
		return typed.TypeArgs()
	case *types.Basic:
		return nil
	default:
		log.Fatalf("%T: %#v", t, t)
	}
	panic("")
}

func getFirstTypeArg[T types.Type](t types.Type) (T, bool) {
	args := getTypeArgs(t)
	if args == nil || args.Len() != 1 {
		var t T
		return t, false
	}
	v, ok := args.At(0).(T)
	return v, ok
}
