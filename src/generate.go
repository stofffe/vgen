package main

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

var tmpl *template.Template

func generateFile(info ParseInfo) ([]byte, error) {
	var buffer bytes.Buffer

	// init global tmpl and parse template
	var err error
	tmpl, err = template.New("template").ParseFiles("src/template.tmpl", "src/rules.tmpl")
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

func templateToString(name string, data any) (string, error) {
	var buffer bytes.Buffer
	err := tmpl.ExecuteTemplate(&buffer, name, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (f PrimField) ValidationCode() (string, error) {
	return templateToString("prim-field-validation", f)
}
func (f TypeField) ValidationCode() (string, error) {
	return templateToString("type-field-validation", f)
}
func (f ListField) ValidationCode() (string, error) {
	return templateToString("list-field-validation", f)
}
func (r Rule) Code() (string, error) {
	return templateToString(r.rule, r)
}

func (f PrimField) ConvertCode() (string, error) {
	return templateToString("prim-field-convert", f)
}
func (f TypeField) ConvertCode() (string, error) {
	return templateToString("type-field-convert", f)
}
func (f ListField) ConvertCode() (string, error) {
	if _, ok := f.inner.(ListField); ok {
		return templateToString("list-field-convert-outer", f)
	}
	return templateToString("list-field-convert-inner", f)
}
