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
var template_str string

//go:embed rules.tmpl
var rules_str string

var tmpl *template.Template

func generateFile(info ParseInfo) ([]byte, error) {
	var buffer bytes.Buffer

	// init global tmpl and parse template
	var err error
	tmpl, err = template.New("template").Funcs(template.FuncMap{
		"iterate": func(count int) []int {
			var i int
			var Items []int
			for i = 0; i < count; i++ {
				Items = append(Items, i)
			}
			return Items
		},
		"call": func(name string, data interface{}) (string, error) {
			var buffer bytes.Buffer
			err := tmpl.ExecuteTemplate(&buffer, name, data)
			if err != nil {
				return "", err
			}
			return buffer.String(), nil
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}).Parse(template_str + rules_str)
	if err != nil {
		return []byte{}, fmt.Errorf("could not parse template file: %v", err)
	}

	// package
	err = tmpl.ExecuteTemplate(&buffer, "package", info)
	if err != nil {
		return []byte{}, fmt.Errorf("could not execute package template: %v", err)
	}

	// structs
	for _, struct_type := range info.StructTypes {
		// type
		err := tmpl.ExecuteTemplate(&buffer, "struct-type", struct_type)
		if err != nil {
			return []byte{}, fmt.Errorf("could not execute struct-type template: %v", err)
		}

		// validation
		err = tmpl.ExecuteTemplate(&buffer, "struct-validation", struct_type)
		if err != nil {
			return []byte{}, fmt.Errorf("could not execute struct-validation: %v", err)
		}

		// conversion
		err = tmpl.ExecuteTemplate(&buffer, "struct-convert", struct_type)
		if err != nil {
			return []byte{}, fmt.Errorf("could not execute struct-convert: %v", err)
		}

		// validation and conversion
		err = tmpl.ExecuteTemplate(&buffer, "struct-validation-convert", struct_type)
		if err != nil {
			return []byte{}, fmt.Errorf("could not execute struct-convert: %v", err)
		}

		// json decoding
		// err = tmpl.ExecuteTemplate(&buffer, "json-decoding", struct_type)
		// if err != nil {
		// 	return []byte{}, fmt.Errorf("could not execute json-decoding: %v", err)
		// }
	}

	// debug
	// return buffer.Bytes(), nil
	if DEBUG {
		return buffer.Bytes(), nil
	}

	// fmt
	bytes, err := format.Source(buffer.Bytes())
	if err != nil {
		return []byte{}, fmt.Errorf("could not format generated file: %v", err)
	}

	// Remove all empty lines
	// TODO slow
	bytes = []byte(strings.ReplaceAll(string(bytes), "\n\n", "\n"))

	return bytes, nil

}
