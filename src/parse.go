package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"
)

type ParseInfo struct {
	Package string
	Structs []StructType
}

type StructType struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name     string
	Type     string
	Rules    []string
	Required bool
}

func parseFile(path string) (ParseInfo, error) {
	// load file
	fset := token.NewFileSet()
	opts := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, path, nil, opts)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not parse file: %v", err)
	}

	err = ast.Print(fset, file)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not print ast: %v", err)
	}

	// parse
	var structs []StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if node, ok := n.(*ast.GenDecl); ok {
			parsed_structs, err := parseGenDecl(node)
			if err != nil {
				fmt.Printf("could not parse gen decl: %v\n", err)
			}
			for _, s := range parsed_structs {
				fmt.Println(s)
				structs = append(structs, s)
			}
		}

		return true
	})

	// package
	package_name := file.Name.Name

	return ParseInfo{
		Package: package_name,
		Structs: structs,
	}, nil
}

const INCLUDE_TAG = `vgen:"i"`

func parseGenDecl(node *ast.GenDecl) ([]StructType, error) {
	// check for tag
	if node.Doc == nil {
		return []StructType{}, nil
	}
	hasTag := false
	for _, comment := range node.Doc.List {
		if strings.Contains(comment.Text, INCLUDE_TAG) {
			hasTag = true
			break
		}
	}
	if !hasTag {
		return []StructType{}, nil
	}

	// check if type
	if node.Tok != token.TYPE {
		return []StructType{}, nil
	}

	// parse all types under decl
	var structs []StructType
	for _, spec := range node.Specs {
		type_node := spec.(*ast.TypeSpec) // already checked

		s, err := parseType(type_node)
		structs = append(structs, s)
		if err != nil {
			return []StructType{}, fmt.Errorf("could not parse type: %v", err)
		}
	}

	return structs, nil
}

func parseType(node *ast.TypeSpec) (StructType, error) {
	// check name
	if node.Name == nil {
		return StructType{}, fmt.Errorf("must have name")
	}

	// check struct
	struct_node, ok := node.Type.(*ast.StructType)
	if !ok {
		return StructType{}, fmt.Errorf("must be struct, got %s", node.Type)
	}

	// check non empty
	if struct_node.Fields == nil || len(struct_node.Fields.List) == 0 {
		return StructType{}, fmt.Errorf("empty structs not supported")
	}

	// parse
	struct_type := StructType{
		Name:   node.Name.Name,
		Fields: []Field{},
	}
	for _, field_node := range struct_node.Fields.List {
		field, err := parseField(field_node)
		if err != nil {
			return StructType{}, fmt.Errorf("could not parse field: %v", err)
		}
		struct_type.Fields = append(struct_type.Fields, field)
	}

	return struct_type, nil
}

func parseField(node *ast.Field) (Field, error) {
	// check name
	if len(node.Names) == 0 {
		return Field{}, fmt.Errorf("field without name not supported")
	}
	if len(node.Names) > 1 {
		return Field{}, fmt.Errorf("field with multiple names not supported")
	}
	name := node.Names[0].Name

	// extract rules from comment
	var rules []string
	if node.Comment != nil {
		if len(node.Comment.List) > 1 {
			log.Fatalf("HOW CAN THIS BE > 1?")
		}

		comment := node.Comment.List[0].Text
		rules = extractRules(comment)
	}

	// parse
	prim, ok := node.Type.(*ast.Ident)
	if !ok {
		return Field{}, fmt.Errorf("must be primitive value")
	}
	field, err := parseFieldPrimitive(prim, name, rules)
	if err != nil {
		return Field{}, fmt.Errorf("could not parse field primtive: %v", err)
	}

	return field, nil
}

func parseFieldPrimitive(node *ast.Ident, field_name string, rules []string) (Field, error) {
	typ := node.Name

	req := false
	for _, rule := range rules {
		if rule == "req" {
			req = true
			break
		}
	}

	return Field{
		Name:     field_name,
		Type:     typ,
		Rules:    rules,
		Required: req,
	}, nil
}

func extractRules(value string) []string {
	reg := regexp.MustCompile(`vgen:\"(.*)\"`)
	matches := reg.FindStringSubmatch(value)

	if len(matches) == 0 {
		return []string{}
	}

	rules := matches[1] // first match is whole string
	rules = strings.ReplaceAll(rules, " ", "")
	split := strings.Split(rules, ",")

	return split
}
