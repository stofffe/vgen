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

func main() {
	err := parseFile("examples/test.go")
	if err != nil {
		log.Fatal(err)
	}
}

func parseFile(path string) error {
	fset := token.NewFileSet()
	opts := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, path, nil, opts)
	if err != nil {
		return fmt.Errorf("could not parse file: %v", err)
	}

	err = ast.Print(fset, file)
	if err != nil {
		return fmt.Errorf("could not print ast: %v", err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if node, ok := n.(*ast.GenDecl); ok {
			err := parseGenDecl(node)
			if err != nil {
				fmt.Printf("could not parse gen decl: %v\n", err)
			}
		}

		return true
	})

	return nil
}

const INCLUDE_TAG = `vgen:"i"`

func parseGenDecl(node *ast.GenDecl) error {
	// check for comment, needed for tag
	if node.Doc == nil {
		return nil
	}

	// check for tag
	hasTag := false
	for _, comment := range node.Doc.List {
		if strings.Contains(comment.Text, INCLUDE_TAG) {
			hasTag = true
			break
		}
	}
	if !hasTag {
		return nil
	}

	// check if type
	if node.Tok != token.TYPE {
		return nil
	}

	// parse all types under decl
	for _, spec := range node.Specs {
		type_node, ok := spec.(*ast.TypeSpec)
		if !ok {
			return fmt.Errorf("spec contains non type which should not be possible") // was checked in prev if
		}

		err := parseType(type_node)
		if err != nil {
			return fmt.Errorf("could not parse type: %v", err)
		}
	}

	return nil
}

func parseType(node *ast.TypeSpec) error {
	fmt.Printf("parse type: %s\n", node.Name)

	struct_node, ok := node.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("currently only supports structs, got %s", node.Type)
	}

	err := parseStruct(struct_node)
	if err != nil {
		return fmt.Errorf("could not parse struct: %v", err)
	}

	return nil
}

func parseStruct(node *ast.StructType) error {
	if node.Fields == nil || len(node.Fields.List) == 0 {
		return fmt.Errorf("empty structs not supported")
	}

	for _, field := range node.Fields.List {
		err := parseField(field)
		if err != nil {
			return fmt.Errorf("could not parse field: %v", err)
		}
	}

	return nil
}

func parseField(node *ast.Field) error {
	// name
	if len(node.Names) == 0 {
		return fmt.Errorf("field without name not supported")
	}
	if len(node.Names) > 1 {
		return fmt.Errorf("field with multiple names not supported")
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

	// parse field
	var err error
	switch n := node.Type.(type) {
	case *ast.Ident:
		err = parseFieldPrimitive(n, name, rules)
	default:
		err = fmt.Errorf("unsupported field type: %T", n)
	}
	if err != nil {
		return err
	}

	return nil
}

func extractRules(value string) []string {
	reg := regexp.MustCompile(`vgen:\[(.*)\]`)
	matches := reg.FindStringSubmatch(value)

	if len(matches) == 0 {
		return []string{}
	}

	rules := matches[1] // first match is whole string
	rules = strings.ReplaceAll(rules, " ", "")
	split := strings.Split(rules, ",")

	return split
}

func parseFieldPrimitive(node *ast.Ident, name string, rules []string) error {
	typ := node.Name
	fmt.Printf("%s %v %v\n", name, typ, rules)

	return nil
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name  string
	Type  string
	Rules []string
}

// func parseFieldStruct(node *ast.StructType) error {
// 	if node.Fields == nil || len(node.Fields.List) == 0 {
// 		return fmt.Errorf("empty structs not supported")
// 	}
//
// 	for _, field := range node.Fields.List {
// 		err := parseField(field)
// 		if err != nil {
// 			return fmt.Errorf("could not parse field: %v", err)
// 		}
// 	}
//
// 	return nil
// }
