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
		parserConfig  = parser.DefaultConfig
		codegenConfig = codegen.DefaultConfig
	)
	pflag.BoolVarP(&format, "format", "f", false, "run go fmt")
	pflag.BoolVarP(&fumpt, "fumpt", "F", false, "run gofumpt")
	pflag.BoolVarP(&imports, "imports", "i", false, "run goimports")
	pflag.BoolVarP(&parserConfig.Verbose, "verbose", "v", false, "verbose logging")
	pflag.Parse()

	if pflag.NArg() > 0 {
		directory = pflag.Arg(0)
	}

	if fumpt && format {
		fmt.Fprintln(os.Stderr, "cannot use go fmt and gofumpt together")
		os.Exit(1)
	}

	if fumpt {
		codegenConfig.OutputFileOptions = append(codegenConfig.OutputFileOptions, codegen.WithGoFumpt())
	}
	if format {
		codegenConfig.OutputFileOptions = append(codegenConfig.OutputFileOptions, codegen.WithGoFmt())
	}
	if imports {
		codegenConfig.OutputFileOptions = append(codegenConfig.OutputFileOptions, codegen.WithGoImports())
	}

	r, err := parserConfig.Parse(directory)
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
				of = codegenConfig.NewOutputFile(s.TargetFile, pkg)
				out[s.TargetFile] = of
				defer of.Close()
				cleanup = append(cleanup, s.TargetFile)
			}
			if parserConfig.Verbose {
				log.Printf("generating %s", s.Name)
			}
			for _, g := range codegenConfig.Generators {
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
