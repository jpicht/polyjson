package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jpicht/polyjson/parser"
	"github.com/jpicht/polyjson/parser/att"
	"github.com/stoewer/go-strcase"
)

type PolyStruct struct {
	Name       string
	TargetFile string
	Package    *att.Package
	Common     *att.NamedStruct
	Additional *att.Field
	Impls      []PolyImpl
}

type PolyImpl struct {
	Struct *att.NamedStruct
	Value  att.MarkerValue
}

func PolyStructs(pkg *att.Package) (out []*PolyStruct, _ error) {
	cache := map[string]*PolyStruct{}
	ps := func(target string) *PolyStruct {
		ps, ok := cache[target]
		if !ok {
			ps = &PolyStruct{
				Name:    target,
				Package: pkg,
			}
			cache[target] = ps
			out = append(out, ps)
		}
		return ps
	}

	for _, s := range pkg.NamedStructs {
		for m, mv := range s.Markers {
			switch m {
			case parser.MarkerIs:
				ps := ps(mv.Target)
				ps.Impls = append(ps.Impls, PolyImpl{
					Struct: s,
					Value:  mv,
				})
			case parser.MarkerCommon:
				ps := ps(mv.Target)
				if ps.Common != nil {
					return nil, fmt.Errorf("two commons for %q: %q and %q", mv, ps.Common, s)
				}
				ps.Common = s
			}
		}
	}

	for _, ps := range out {
		ps.TargetFile = filepath.Join(filepath.Dir(ps.Package.GoFiles[0]), strcase.SnakeCase(ps.Name)+"_polyjson.go")

		if ps.Common == nil {
			continue
		}

		ps.TargetFile = strings.TrimSuffix(ps.Common.File, ".go") + "_polyjson.go"

		for _, f := range ps.Common.Fields {
			if !f.Interfaces[parser.AdditionalFields] {
				continue
			}

			if ps.Additional != nil {
				return nil, fmt.Errorf("two additional field sinks in %q: %q and %q", ps.Name, ps.Additional.Name, f.Name)
			}

			ps.Additional = f
		}
	}

	return
}