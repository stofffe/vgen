package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

const includeTag = "vgen"

type ParseInfo struct {
	Package     string
	StructTypes []StructType
	Imports     []string
}

type StructType struct {
	Name   string
	Fields []StructField
}

type StructField struct {
	Name   string
	Type   string
	Tags   string
	Alias  string // json tag
	Nested bool
}

func parseFile(path string) (ParseInfo, error) {
	// load file
	fset := token.NewFileSet()
	opts := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, path, nil, opts)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not parse file: %v", err)
	}

	// traverse tree
	var fileErr error
	var structTypes []StructType
	ast.Inspect(file, func(n ast.Node) bool {
		node, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		if node.Tok == token.TYPE {
			// Check for tag
			comment := node.Doc.Text()
			if !strings.Contains(comment, includeTag) {
				return true
			}

			// Parse type
			parsedTypes, err := parseType(node)
			if err != nil {
				fileErr = fmt.Errorf("could not parse type: %v", err)
				return false
			}
			for _, s := range parsedTypes {
				structTypes = append(structTypes, s)
			}
		}

		if node.Tok == token.IMPORT {
			// fmt.Println("imports", node)
		}

		return true
	})
	if fileErr != nil {
		return ParseInfo{}, fileErr
	}

	packageName := file.Name.Name

	return ParseInfo{
		Package:     packageName,
		StructTypes: structTypes,
		Imports: []string{
			"github.com/stofffe/vgen",
		},
	}, nil
}

func parseType(declNode *ast.GenDecl) ([]StructType, error) {
	var structs []StructType
	for _, spec := range declNode.Specs {
		typeNode := spec.(*ast.TypeSpec)

		// check name
		if typeNode.Name == nil {
			return []StructType{}, fmt.Errorf("must have name")
		}

		structNode, ok := typeNode.Type.(*ast.StructType)
		if !ok {
			return []StructType{}, fmt.Errorf("unsupported type")
		}

		name := typeNode.Name.Name
		structType, err := parseStruct(structNode, name)
		if err != nil {
			return []StructType{}, fmt.Errorf("could not parse struct: %v", err)
		}

		structs = append(structs, structType)
	}

	return structs, nil
}

func parseStruct(structNode *ast.StructType, structName string) (StructType, error) {
	structType := StructType{
		Name:   structName,
		Fields: []StructField{},
	}

	for _, fieldNode := range structNode.Fields.List {
		field, err := parseField(fieldNode)
		if err != nil {
			return StructType{}, fmt.Errorf("could not parse field: %v", err)
		}
		structType.Fields = append(structType.Fields, field)
	}

	return structType, nil
}

func parseField(fieldNode *ast.Field) (StructField, error) {
	comments := fieldNode.Doc.Text() + fieldNode.Comment.Text()

	fieldName := fieldNode.Names[0].Name // TODO handle multiple
	tags := ""
	if fieldNode.Tag != nil {
		tags = fieldNode.Tag.Value
	}

	commentTags, err := parseFieldTags(comments)
	if err != nil {
		return StructField{}, fmt.Errorf("could not parse field tags: %v", err)
	}
	nested := commentTags.include
	alias := fieldName
	if commentTags.name != "" {
		alias = commentTags.name
	}

	typ, err := parseFieldType(fieldNode.Type)
	if err != nil {
		return StructField{}, fmt.Errorf("could not parse field type: %v", err)
	}

	return StructField{
		Name:   fieldName,
		Type:   typ,
		Tags:   tags,
		Alias:  alias,
		Nested: nested,
	}, nil
}

type FieldTags struct {
	include bool
	name    string
}

func parseFieldTags(comment string) (FieldTags, error) {
	// default tags
	tags := FieldTags{
		name:    "",
		include: false,
	}

	reg := regexp.MustCompile(`vgen\((?s).*\)`)
	match := reg.FindString(comment)

	if match == "" {
		return tags, nil
	}

	match = strings.TrimPrefix(match, "vgen(")
	match = strings.TrimSuffix(match, ")")
	args := strings.Split(match, ",")

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		split := strings.Split(arg, "=")
		ident := split[0]
		switch ident {
		case "nested", "n":
			tags.include = true
		case "alias":
			if len(split) < 2 || split[1] == "" {
				return FieldTags{}, fmt.Errorf("name must have second argument")
			}
			name := split[1]
			tags.name = name
		default:
			return FieldTags{}, fmt.Errorf("unknown tag %v", ident)
		}
	}
	return tags, nil
}

func parseFieldType(fieldNode ast.Expr) (string, error) {
	switch node := fieldNode.(type) {
	case *ast.Ident:
		return node.Name, nil
	case *ast.ArrayType:
		fieldType, err := parseFieldType(node.Elt)
		if err != nil {
			return "", err
		}
		return "[]" + fieldType, nil
	case *ast.MapType:
		key, ok := node.Key.(*ast.Ident)
		if !ok || key.Name != "string" {
			return "", fmt.Errorf("invalid map key %v, must be string", key)
		}

		fieldType, err := parseFieldType(node.Value)
		if err != nil {
			return "", err
		}
		return "map[string]" + fieldType, nil
	case *ast.StarExpr:
		fieldType, err := parseFieldType(node.X)
		if err != nil {
			return "", err
		}
		return "*" + fieldType, nil
	}
	return "", fmt.Errorf("unsupported field type: %T", fieldNode)
}

func extractJsonName(tag string) (string, bool) {
	reg := regexp.MustCompile(`json:"[^"]*"`)
	match := reg.FindString(tag)
	match = strings.TrimPrefix(match, `json:"`)
	match = strings.TrimSuffix(match, `"`)
	match = strings.Split(match, ",")[0]
	ok := match != ""
	return match, ok
}
