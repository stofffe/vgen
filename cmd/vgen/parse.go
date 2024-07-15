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
	Name    string
	Types   []FieldType
	Tags    string
	Alias   string
	Nested  bool
	Pointer bool
}

func (s StructField) Type() string {
	var builder strings.Builder
	if s.Pointer {
		builder.WriteByte('*')
	}
	for _, t := range s.Types {
		builder.WriteString(t.Type())
	}
	return builder.String()
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
			typeTags, err := parseTypeTags(comment)
			if err != nil {
				fileErr = fmt.Errorf("could not parse type tags: %v", err)
				return false
			}
			if !typeTags.include {
				return false
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

	// parse tags
	commentTags, err := parseFieldTags(comments)
	if err != nil {
		return StructField{}, fmt.Errorf("could not parse field tags: %v", err)
	}
	nested := commentTags.include
	alias := fieldName
	if commentTags.alias != "" {
		alias = commentTags.alias
	}

	// parse field types
	fieldTypeInfo := FieldTypeInfo{
		Types:     []FieldType{},
		Pointer:   false,
		Primitive: false,
	}
	err = fieldTypeInfo.parseFieldType(fieldNode.Type)
	if err != nil {
		return StructField{}, fmt.Errorf("could not parse field type: %v", err)
	}

	// Dont allow nested on primitve types
	if nested && fieldTypeInfo.Primitive {
		return StructField{}, fmt.Errorf("nested not allowed on primitve inner type")
	}

	return StructField{
		Name:   fieldName,
		Tags:   tags,
		Alias:  alias,
		Nested: nested,

		Types:   fieldTypeInfo.Types,
		Pointer: fieldTypeInfo.Pointer,
	}, nil
}

type FieldType interface {
	Type() string
}

type FieldTypeIdent struct{ name string }
type FieldTypeArray struct{}
type FieldTypeMap struct{}

func (f FieldTypeIdent) Type() string { return f.name }
func (f FieldTypeArray) Type() string { return "[]" }
func (f FieldTypeMap) Type() string   { return "map[string]" }

type FieldTypeInfo struct {
	Types     []FieldType
	Pointer   bool
	Primitive bool
}

func (f *FieldTypeInfo) parseFieldType(fieldNode ast.Expr) error {
	switch node := fieldNode.(type) {
	case *ast.Ident:
		if node.Obj == nil {
			f.Primitive = true
		}
		f.Types = append(f.Types, FieldTypeIdent{name: node.Name})
		return nil
	case *ast.ArrayType:
		f.Types = append(f.Types, FieldTypeArray{})
		err := f.parseFieldType(node.Elt)
		if err != nil {
			return err
		}
	case *ast.MapType:
		f.Types = append(f.Types, FieldTypeMap{})
		key, ok := node.Key.(*ast.Ident)
		if !ok || key.Name != "string" {
			return fmt.Errorf("invalid map key %v, must be string", key)
		}

		err := f.parseFieldType(node.Value)
		if err != nil {
			return err
		}
	case *ast.StarExpr:
		if len(f.Types) > 0 {
			return fmt.Errorf("pointers not allowed as list or map element")
		}
		err := f.parseFieldType(node.X)
		if err != nil {
			return err
		}
		f.Pointer = true
	default:
		return fmt.Errorf("unsupported field type: %T", fieldNode)
	}
	return nil
}

type FieldTags struct {
	include bool
	alias   string
}

func parseFieldTags(comment string) (FieldTags, error) {
	// default tags
	tags := FieldTags{
		alias:   "",
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
			tags.alias = name
		default:
			return FieldTags{}, fmt.Errorf("unknown tag %v", ident)
		}
	}
	return tags, nil
}

type TypeTags struct {
	include bool
}

func parseTypeTags(comment string) (TypeTags, error) {
	// default tags
	tags := TypeTags{
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
		case "include", "i":
			tags.include = true
		default:
			return TypeTags{}, fmt.Errorf("unknown type tag %v", ident)
		}
	}
	return tags, nil
}
