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
		"RuleFunc":       ruleFunc,
		"lowerFirstFunc": lowerFirstFunc,
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

	// fmt
	bytes, err := format.Source(buffer.Bytes())
	if err != nil {
		return []byte{}, fmt.Errorf("could not format generated file: %v", err)
	}

	return bytes, nil
}

// TODO temp
func ruleFunc(rule string) string {
	if rule == "req" {
		return ""
	}
	if rule == "len>0" {
		return `
		if !(len(name) > 0) {
			errs[""] = fmt.Sprintf("len must be > 0")
		}`
	}

	return "// Rule not implemented for " + rule
}

func lowerFirstFunc(str string) string {
	if str == "" {
		return str
	}
	firstchar := []rune(str)[0]
	firstchar = unicode.ToLower(firstchar)
	return string(firstchar) + str[1:]
}
