package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

func (r *Result) parseFile(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	if additional, ok := r.Additional[filename]; ok {
		src = append(src, additional...)
	}

	if strings.HasSuffix(filename, "_polyjson.go") || strings.HasSuffix(filename, "_mock.go") {
		log.Printf("pre-parsing file %q", filename)
		scanner := bufio.NewScanner(bytes.NewBuffer(src))
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.HasPrefix(txt, "package ") {
				src = []byte(txt + "\n")
				break
			}
		}
	}

	f, err := parser.ParseFile(fset, filename, src, parser.AllErrors)
	return f, err
}

func (r *Result) load() (err error) {
	r.Config.Packages.Fset = r.StdLib.FileSet
	r.Config.Packages.ParseFile = r.parseFile

	for {
		if r.Config.Verbose {
			log.Println("load")
		}
		r.StdLib.Packages, err = packages.Load(&r.Config.Packages)

		missing := make(map[string]string)
		for _, inPkg := range r.StdLib.Packages {
			if r.Config.Verbose {
				log.Printf("  pkg %q", inPkg.Name)
			}

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
			if r.Config.Verbose {
				log.Printf("adding %q to %q", name, file)
			}
			r.Additional[file] = append(r.Additional[file], []byte(
				fmt.Sprintf("\n\ntype %s struct {}", name),
			)...)
		}

		if len(missing) == 0 {
			return
		}
	}
}
