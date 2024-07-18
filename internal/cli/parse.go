package cli

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
	Imports     map[string]string
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
	Import  string
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
		return ParseInfo{}, fmt.Errorf("parse file: %w", err)
	}

	var traverseErr error
	var structTypes []StructType
	importMap := make(map[string]string)

	// traverse tree
	ast.Inspect(file, func(n ast.Node) bool {
		node, ok := n.(*ast.GenDecl)
		if !ok {
			return true
		}

		// type
		if node.Tok == token.TYPE {
			// Check for tag
			comment := node.Doc.Text()
			typeTags, err := parseTypeTags(comment)
			if err != nil {
				traverseErr = fmt.Errorf("parse type tags: %w", err)
				return false
			}
			if !typeTags.include {
				return false
			}

			// Parse type
			parsedTypes, err := parseType(node)
			if err != nil {
				traverseErr = fmt.Errorf("parse type: %w", err)
				return false
			}
			for _, s := range parsedTypes {
				structTypes = append(structTypes, s)
			}
		}

		// imports
		if node.Tok == token.IMPORT {
			for _, spec := range node.Specs {
				importSpec, _ := spec.(*ast.ImportSpec)
				if !ok {
					traverseErr = fmt.Errorf("invalid import spec %T", spec)
					return false
				}

				path := importSpec.Path.Value
				path = strings.TrimPrefix(path, `"`)
				path = strings.TrimSuffix(path, `"`)
				name := ""
				if importSpec.Name != nil {
					// custom name
					name = importSpec.Name.Name
				} else {
					// strip last
					split := strings.Split(path, "/")
					name = split[len(split)-1]
				}

				// remove quotes
				importMap[name] = path
			}
		}

		return true
	})
	if traverseErr != nil {
		return ParseInfo{}, traverseErr
	}

	// Save all imports used in vgen type fields
	imports := map[string]string{
		"vgen": "github.com/stofffe/vgen/pkg/vgen",
	}
	for _, typ := range structTypes {
		for _, field := range typ.Fields {
			if field.Import == "" {
				continue
			}
			imports[field.Import] = importMap[field.Import]
		}
	}

	packageName := file.Name.Name
	return ParseInfo{
		Package:     packageName,
		StructTypes: structTypes,
		Imports:     imports,
	}, nil
}

func parseType(declNode *ast.GenDecl) ([]StructType, error) {
	var structs []StructType
	for _, spec := range declNode.Specs {
		typeNode := spec.(*ast.TypeSpec)

		// check name
		if typeNode.Name == nil {
			return []StructType{}, DetailedError{
				msg: "parsed types must have a name",
				err: fmt.Errorf("must have name"),
			}
		}

		switch node := typeNode.Type.(type) {
		case *ast.StructType:
			name := typeNode.Name.Name
			structType, err := parseStruct(node, name)
			if err != nil {
				return []StructType{}, fmt.Errorf("prase struct: %w", err)
			}
			structs = append(structs, structType)
		case *ast.Ident:
			return nil, DetailedError{
				msg: "type aliases not supported",
				err: fmt.Errorf("type aliases not supported"),
			}
		default:
			return nil, fmt.Errorf("unsupported type %T", node)
		}
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
			return StructType{}, fmt.Errorf("parse field: %w", err)
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
		return StructField{}, fmt.Errorf("parse field tags: %w", err)
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
		Import:    "",
	}
	err = fieldTypeInfo.parseFieldType(fieldNode.Type)
	if err != nil {
		return StructField{}, fmt.Errorf("field type: %w", err)
	}

	// Dont allow nested on primitve types
	if nested && fieldTypeInfo.Primitive {
		return StructField{}, DetailedError{
			msg: "primitve fields can not have nested tag",
			err: fmt.Errorf("nested not allowed on primitve inner type"),
		}
	}

	return StructField{
		Name:   fieldName,
		Tags:   tags,
		Alias:  alias,
		Nested: nested,

		Types:   fieldTypeInfo.Types,
		Pointer: fieldTypeInfo.Pointer,
		Import:  fieldTypeInfo.Import,
	}, nil
}

type FieldType interface {
	Type() string
}

type FieldTypeIdent struct{ name string }
type FieldTypeImport struct {
	imp  string
	name string
}
type FieldTypeArray struct{}
type FieldTypeMap struct{}

func (f FieldTypeIdent) Type() string  { return f.name }
func (f FieldTypeImport) Type() string { return f.imp + "." + f.name }
func (f FieldTypeArray) Type() string  { return "[]" }
func (f FieldTypeMap) Type() string    { return "map[string]" }

type FieldTypeInfo struct {
	Types     []FieldType
	Pointer   bool
	Primitive bool
	Import    string
}

func (f *FieldTypeInfo) parseFieldType(fieldNode ast.Expr) error {
	switch node := fieldNode.(type) {
	// internal primitive/struct
	case *ast.Ident:
		if node.Obj == nil {
			f.Primitive = true
		}
		f.Types = append(f.Types, FieldTypeIdent{name: node.Name})
	// array
	case *ast.ArrayType:
		f.Types = append(f.Types, FieldTypeArray{})
		err := f.parseFieldType(node.Elt)
		if err != nil {
			return fmt.Errorf("array: %w", err)
		}
	// map
	case *ast.MapType:
		f.Types = append(f.Types, FieldTypeMap{})
		key, ok := node.Key.(*ast.Ident)
		if !ok || key.Name != "string" {
			return DetailedError{
				msg: "key of map must be a string",
				err: fmt.Errorf("invalid map key %v, must be string", key),
			}
		}

		err := f.parseFieldType(node.Value)
		if err != nil {
			return fmt.Errorf("map field: %w", err)
		}
	// pointer
	case *ast.StarExpr:
		if len(f.Types) > 0 {
			return DetailedError{
				msg: "pointers not allowed as list/map element",
				err: fmt.Errorf("pointers not allowed as list/map element"),
			}
		}
		err := f.parseFieldType(node.X)
		if err != nil {
			return err
		}
		f.Pointer = true
	// import
	case *ast.SelectorExpr:
		imp, ok := node.X.(*ast.Ident)
		if !ok {
			return fmt.Errorf("import selector is not ast.Ident")
		}
		f.Import = imp.Name
		f.Types = append(f.Types, FieldTypeImport{imp: imp.Name, name: node.Sel.Name})
	default:
		return fmt.Errorf("unknown field type %T", fieldNode)
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
				return FieldTags{}, DetailedError{
					msg: `field tag alias must have a value, ex: "alias=something"`,
					err: fmt.Errorf("alias must have second argument"),
				}
			}
			name := split[1]
			tags.alias = name
		default:
			return FieldTags{}, DetailedError{
				msg: fmt.Sprintf("unknown field tag %v", ident),
				err: fmt.Errorf("unknown field tag %v", ident),
			}
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
			return TypeTags{}, DetailedError{
				msg: fmt.Sprintf("unknown type tag %s", ident),
				err: fmt.Errorf("unknown type tag %s", ident),
			}
		}
	}
	return tags, nil
}
