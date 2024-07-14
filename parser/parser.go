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

const debug = false

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
		InterfaceTests  map[att.Marker]func(types.Type) bool
		TypeMarkerTests map[att.Marker]func(types.Type) bool
		MarkerTests     map[att.Marker]func(types.Type) bool
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

	if err := r.load(); err != nil {
		return nil, err
	}

	r.Helpers.MarkerTests, _ = resolveMarkers(config, config.Markers, r.StdLib.Packages)
	r.Helpers.TypeMarkerTests, _ = resolveMarkers(config, config.TypeMarkers, r.StdLib.Packages)
	r.Helpers.InterfaceTests, _ = resolveMarkers(config, config.Interfaces, r.StdLib.Packages)

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

		isMarker := false
		for marker, tm := range r.Helpers.TypeMarkerTests {
			if !tm(obj.Type()) {
				continue
			}

			tt := obj.Type()

			t, ok := getFirstTypeArg[*types.Named](tt)
			if !ok {
				continue
			}

			outPkg.TypeMarkers = append(outPkg.TypeMarkers, att.TypeMarker{
				Marker: marker,
				Target: t.Obj().Name(),
				Name:   name,
				Type:   tt,
			})

			isMarker = true

			log.Printf("type marker %s: %s for %s ", marker.InterfaceName, name, t.Obj().Name())
		}

		if isMarker {
			continue
		}

		if debug {
			log.Printf("obj[%s]: %T -> %T", types.TypeString(obj.Type(), qual), obj, obj.Type())
			log.Printf("  %s", types.TypeString(obj.Type().Underlying(), qual))
		}

		named, structT, ok := getStruct(obj.Type())
		if !ok {
			if r.Config.Verbose {
				log.Printf("skipping %q (%s)", name, types.TypeString(obj.Type(), qual))
			}
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
						if r.Config.Verbose {
							log.Printf("skipping field %s on %s, no type arg", f.Name, outStruct.Name)
						}
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
				if r.Config.Verbose {
					log.Printf("  marker %s on %s", f.Name, outStruct.Name)
				}
				continue
			}

			for iface, testFn := range r.Helpers.InterfaceTests {
				if !testFn(f.Type) {
					continue
				}
				f.Interfaces[iface] = true
			}

			outStruct.Fields = append(outStruct.Fields, f)

			if r.Config.Verbose {
				log.Printf("  field %s.%s (%s)", outStruct.Name, f.Name, types.TypeString(f.Type, qual))
			}
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

func getNamed(t types.Type) (named *types.Named, ok bool) {
	named, ok = t.(*types.Named)
	return
}

func getStruct(t types.Type) (named *types.Named, structT *types.Struct, ok bool) {
	named, ok = getNamed(t)
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
	case *types.Alias:
		return getTypeArgs(typed.Obj().Type())
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
