package cli

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
	Name string
	Type string
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
			fmt.Println("imports", node)
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
			"github.com/stofffe/vgen/vgen",
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

	for _, field := range structNode.Fields.List {
		fieldName := field.Names[0].Name // TODO handle multiple
		typ, err := parseFieldType(field.Type)
		if err != nil {
			return StructType{}, fmt.Errorf("could not parse field type: %v", err)
		}
		structType.Fields = append(structType.Fields, StructField{
			Name: fieldName,
			Type: typ,
		})
	}

	return structType, nil
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
	}
	return "", fmt.Errorf("unsupported field type")
}
