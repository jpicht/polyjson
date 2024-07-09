package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

func (r *Result) parseFile(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	if additional, ok := r.Additional[filename]; ok {
		src = append(src, additional...)
	}

	f, err := parser.ParseFile(fset, filename, src, parser.AllErrors)
	if err != nil || !strings.HasSuffix(filename, "_polyjson.go") {
		return f, err
	}

	log.Printf("pre-parsing file %q", filename)

	f.Decls = slices.DeleteFunc(f.Decls, func(decl ast.Decl) bool {
		switch typeDecl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range typeDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structSpec, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				structSpec.Fields = nil
				log.Printf("\tremoving fields from struct %v", typeSpec.Name)
			}

		case *ast.BadDecl:
			return true

		case *ast.FuncDecl:
			return true

		default:
		}
		return false
	})

	return f, nil
}

func (r *Result) load() (err error) {
	r.Config.Packages.Fset = r.StdLib.FileSet
	r.Config.Packages.ParseFile = r.parseFile

	for {
		r.StdLib.Packages, err = packages.Load(&r.Config.Packages)

		missing := make(map[string]string)
		for _, inPkg := range r.StdLib.Packages {
			if len(inPkg.TypeErrors) > 0 {
				for _, te := range inPkg.TypeErrors {
					if name, ok := strings.CutPrefix(te.Msg, "undefined: "); ok {
						missing[name] = r.StdLib.FileSet.File(te.Pos).Name()
						continue
					}
					log.Printf("type error %q", te.Msg)
				}
			}
		}

		for name, file := range missing {
			log.Printf("adding %q to %q", name, file)
			r.Additional[file] = append(r.Additional[file], []byte(
				fmt.Sprintf("\n\ntype %s struct {}", name),
			)...)
		}

		if len(missing) == 0 {
			return
		}
	}
}
