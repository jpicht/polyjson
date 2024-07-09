package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/jpicht/polyjson/codegen"
	"github.com/jpicht/polyjson/generator"
	"github.com/jpicht/polyjson/parser"
	"github.com/ogier/pflag"
)

func main() {
	var (
		directory     string
		format        bool
		fumpt         bool
		imports       bool
		outputOptions []codegen.OutputFileOption
	)
	pflag.BoolVarP(&format, "format", "f", false, "run go fmt")
	pflag.BoolVarP(&fumpt, "fumpt", "F", false, "run gofumpt")
	pflag.BoolVarP(&imports, "imports", "i", false, "run goimports")
	pflag.Parse()

	if pflag.NArg() > 0 {
		directory = pflag.Arg(0)
	}

	if fumpt && format {
		fmt.Fprintln(os.Stderr, "cannot use go fmt and gofumpt together")
		os.Exit(1)
	}

	if fumpt {
		outputOptions = append(outputOptions, codegen.WithGoFumpt())
	}
	if format {
		outputOptions = append(outputOptions, codegen.WithGoFmt())
	}
	if imports {
		outputOptions = append(outputOptions, codegen.WithGoImports())
	}

	r, err := parser.DefaultConfig.Parse(directory)
	if err != nil {
		log.Fatal(err)
	}

	var cleanup []string

	defer func() {
		e := recover()
		if e == nil {
			return
		}
		log.Fatal(e, string(debug.Stack()))
		for _, f := range cleanup {
			os.Remove(f)
		}
	}()

	var out = make(map[string]*codegen.OutputFile)
	for _, pkg := range r.Packages {
		poly, err := generator.PolyStructs(pkg)
		if err != nil {
			log.Fatal(err)
		}
		for _, s := range poly {
			of, ok := out[s.TargetFile]
			if !ok {
				of = codegen.NewOutputFile(s.TargetFile, pkg, outputOptions...)
				out[s.TargetFile] = of
				defer of.Close()
				cleanup = append(cleanup, s.TargetFile)
			}
			log.Printf("generating %s", s.Name)
			for _, g := range codegen.All {
				func() {
					defer func() {
						if e := recover(); e != nil {
							if err, ok := e.(error); ok {
								of.AddError(fmt.Errorf("recovered: %w", err))
							} else {
								of.AddError(fmt.Errorf("recovered: %v", e))
							}
						}
					}()
					of.AddError(g.GeneratePolyStruct(of.Context, s))
				}()
			}
		}
	}
}
