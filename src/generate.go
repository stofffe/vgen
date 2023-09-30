package main

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"
	"unicode"
)

func generateFile(info ParseInfo) ([]byte, error) {
	var buffer bytes.Buffer

	func_map := template.FuncMap{
		// "ruleFunc":       ruleFunc,
		"lowerFirstFunc": lowerFirstFunc,
		"tmpl":           Tmpl,
	}

	// parse template
	tmpl, err := template.New("template").Funcs(func_map).ParseFiles("src/template.tmpl")
	// tmpl, err := template.ParseFiles("src/template.tmpl")
	if err != nil {
		return []byte{}, fmt.Errorf("could not parse template file: %v", err)
	}

	// package
	err = tmpl.ExecuteTemplate(&buffer, "package", info.Package)
	if err != nil {
		return []byte{}, fmt.Errorf("could not execute package template: %v", err)
	}

	// structs
	for _, s := range info.Structs {
		// type
		err := tmpl.ExecuteTemplate(&buffer, "struct_type", s)
		if err != nil {
			return []byte{}, fmt.Errorf("could not execute struct_type template: %v", err)
		}

		// validation
		err = tmpl.ExecuteTemplate(&buffer, "struct_validation", s)
		if err != nil {
			return []byte{}, fmt.Errorf("could not execute struct_validation: %v", err)
		}
	}

	// return buffer.Bytes(), nil

	// fmt
	bytes, err := format.Source(buffer.Bytes())
	if err != nil {
		return []byte{}, fmt.Errorf("could not format generated file: %v", err)
	}

	return bytes, nil
}

func (f ListField) RulesCode() (string, error) {
	func_map := template.FuncMap{
		"lowerFirstFunc": lowerFirstFunc,
		"tmpl":           Tmpl,
	}

	tmpl, err := template.New("rules").Funcs(func_map).ParseFiles("src/template.tmpl")
	if err != nil {
		return "", fmt.Errorf("could not parse rules template file: %v", err)
	}

	var buffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&buffer, "list_field_validation", f)
	if err != nil {
		return "", fmt.Errorf("could not execute rules template file: %v", err)
	}

	return buffer.String(), nil

}

func (f PrimitiveField) RulesCode() (string, error) {
	func_map := template.FuncMap{
		"lowerFirstFunc": lowerFirstFunc,
		"tmpl":           Tmpl,
	}

	tmpl, err := template.New("rules").Funcs(func_map).ParseFiles("src/template.tmpl")
	if err != nil {
		return "", fmt.Errorf("could not parse rules template file: %v", err)
	}

	var buffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&buffer, "primitive_field_validation", f)
	if err != nil {
		return "", fmt.Errorf("could not execute rules template file: %v", err)
	}

	return buffer.String(), nil
}

// TODO SLOW reuse template?
// func (rule Rule) RulesCodePrefix(prefix string) (string, error) {
// 	rule.ErrorPrefix = prefix
// 	return rule.RulesCode()
// }

func Tmpl(name string, data any) (string, error) {
	func_map := template.FuncMap{
		"lowerFirstFunc": lowerFirstFunc,
		"tmpl":           Tmpl,
	}

	tmpl, err := template.New("rules").Funcs(func_map).ParseFiles("src/rules.tmpl")
	if err != nil {
		return "", fmt.Errorf("could not parse rules template file: %v", err)
	}

	var buffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&buffer, name, data)
	if err != nil {
		return "", fmt.Errorf("could not execute rules template file: %v", err)
	}

	return buffer.String(), nil
}

func lowerFirstFunc(str string) string {
	if str == "" {
		return str
	}
	firstchar := []rune(str)[0]
	firstchar = unicode.ToLower(firstchar)
	return string(firstchar) + str[1:]
}
