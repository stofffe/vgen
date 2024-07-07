package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

//go:embed template.tmpl
var templateStr string

func GenerateFile(info ParseInfo) ([]byte, error) {

	var buffer bytes.Buffer

	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %v", err)
	}

	// package
	err = tmpl.ExecuteTemplate(&buffer, "package", info)
	if err != nil {
		return nil, fmt.Errorf("could not execute template package: %v", err)
	}

	for _, structType := range info.StructTypes {
		// struct type
		err = tmpl.ExecuteTemplate(&buffer, "structType", structType)
		if err != nil {
			return nil, fmt.Errorf("could not execute template structType: %v", err)
		}

		// struct validation
		err = tmpl.ExecuteTemplate(&buffer, "structValidation", structType)
		if err != nil {
			return nil, fmt.Errorf("could not execute template structValidation: %v", err)
		}

		// rule type
		err = tmpl.ExecuteTemplate(&buffer, "ruleType", structType)
		if err != nil {
			return nil, fmt.Errorf("could not execute template ruleType: %v", err)
		}
	}

	bytes := buffer.Bytes()
	// fmt.Println(string(bytes))

	// Format
	bytes, err = format.Source(buffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("could not format generated file")
	}

	// Remove all empty lines
	bytes = []byte(strings.ReplaceAll(string(bytes), "\n\n", "\n"))

	return bytes, nil
}
