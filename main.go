package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputName    string
	packageName   string
	componentName string
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.StringVar(&outputName, "o", "", "output filename")
	flag.StringVar(&packageName, "p", "main", "output package name")
	flag.StringVar(&componentName, "c", "", "component name")
	flag.Parse()
	inputName := flag.Arg(0)
	baseName := inputName[:len(inputName)-len(filepath.Ext(inputName))]
	name := filepath.Base(baseName)
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
			outputName = baseName + "_gen.go"
		} else {
			outputName = "generated.go"
		}
	}
	if len(componentName) == 0 {
		componentName = strings.Title(name)
	}
	output, err := os.Create(outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	log.Printf("gen: %s -> %s", inputName, outputName)
	converter := New()
	buffer := bytes.NewBuffer(nil)
	if err := converter.Do(buffer, input, packageName); err != nil {
		log.Fatal(err)
	}
	if err := templ.Execute(output, map[string]interface{}{
		"PkgName":       packageName,
		"StdImports":    converter.StdModules,
		"Imports":       converter.ExtModules,
		"ComponentName": componentName,
		"Generated":     buffer.String(),
		"Methods":       converter.Methods,
	}); err != nil {
		log.Fatal(err)
	}
}
