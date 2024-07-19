package cli

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

func generateFile(info ParseInfo) ([]byte, error) {

	var buffer bytes.Buffer

	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"iter": func(count int) []int {
			var Items []int
			for i := 0; i < count; i++ {
				Items = append(Items, i)
			}
			return Items
		},
		"iterRange": func(start, end int) []int {
			var Items []int
			for i := start; i < end; i++ {
				Items = append(Items, i)
			}
			return Items
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}).Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("parse template: %v", err)
	}

	// package
	err = tmpl.ExecuteTemplate(&buffer, "package", info)
	if err != nil {
		return nil, fmt.Errorf("execute template package: %v", err)
	}

	for _, structType := range info.StructTypes {
		// struct type
		err = tmpl.ExecuteTemplate(&buffer, "structType", structType)
		if err != nil {
			return nil, fmt.Errorf("execute template structType: %v", err)
		}

		// validation
		err = tmpl.ExecuteTemplate(&buffer, "validation", structType)
		if err != nil {
			return nil, fmt.Errorf("execute template validation: %v", err)
		}

		// rule type
		err = tmpl.ExecuteTemplate(&buffer, "ruleType", structType)
		if err != nil {
			return nil, fmt.Errorf("execute template ruleType: %v", err)
		}

		// convert
		err = tmpl.ExecuteTemplate(&buffer, "convert", structType)
		if err != nil {
			return nil, fmt.Errorf("execute template convert: %v", err)
		}

		// validated convert
		err = tmpl.ExecuteTemplate(&buffer, "validatedConvert", structType)
		if err != nil {
			return nil, fmt.Errorf("execute template validatedConvert: %v", err)
		}
	}

	bytes := buffer.Bytes()

	// Format
	bytes, err = format.Source(buffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format generated file")
	}

	// Remove all empty lines
	bytes = []byte(strings.ReplaceAll(string(bytes), "\n\n", "\n"))

	return bytes, nil
}
