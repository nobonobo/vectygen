package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	outputName    string
	packageName   string
	componentName string
)

func do(output io.Writer, input io.Reader, pkg string) error {
	if err := generate(output, input); err != nil && err != io.EOF {
		return err
	}
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.StringVar(&outputName, "o", "", "output filename")
	flag.StringVar(&packageName, "p", "main", "output package name")
	flag.StringVar(&componentName, "c", "Component", "component name")
	flag.Parse()
	inputName := flag.Arg(0)
	var input io.Reader = os.Stdin
	if len(inputName) > 0 && inputName != "-" {
		r, err := os.Open(inputName)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		input = r
	}
	if len(outputName) == 0 {
		if len(inputName) > 0 {
			outputName = inputName[:len(inputName)-len(filepath.Ext(inputName))] + "_gen.go"
		} else {
			outputName = "generated.go"
		}
	}
	output, err := os.Create(outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()
	log.Printf("gen: %s -> %s", inputName, outputName)
	buffer := bytes.NewBuffer(nil)
	if err := do(buffer, input, packageName); err != nil {
		log.Fatal(err)
	}
	stdImports := []string{}
	imports := []string{
		"github.com/gopherjs/vecty",
		"github.com/gopherjs/vecty/elem",
		"github.com/gopherjs/vecty/prop",
	}
	methods := map[string]struct{}{
		"OnClick":    struct{}{},
		"OnDblClick": struct{}{},
	}
	if len(methods) > 0 {
		stdImports = append(stdImports, "log")
	}
	if err := templ.Execute(output, map[string]interface{}{
		"PkgName":       packageName,
		"StdImports":    stdImports,
		"Imports":       imports,
		"ComponentName": componentName,
		"Generated":     buffer.String(),
		"Methods":       methods,
	}); err != nil {
		log.Fatal(err)
	}
}
