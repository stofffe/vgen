package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/stofffe/vgen/util"
)

type ParseInfo struct {
	Package     string
	StructTypes []StructType
	Imports     []string
}

type StructType struct {
	name   string
	fields []StructField
}

func (f StructType) Name() string          { return f.name }
func (f StructType) LowerName() string     { return util.LowerFirstChar(f.name) }
func (f StructType) Fields() []StructField { return f.fields }

type StructField struct {
	name     string
	required bool
	field    Field
}

func (f StructField) Name() string      { return f.name }
func (f StructField) LowerName() string { return util.LowerFirstChar(f.name) }
func (f StructField) Field() Field      { return f.field }
func (f StructField) Required() bool    { return f.required }

type Field interface {
	Typ() string
	Code() (string, error)
}

type PrimField struct {
	name  string
	typ   string
	rules []Rule
}

type ListField struct {
	name  string
	inner Field
	depth int
	rules []Rule
}

func (f ListField) Inner() Field {
	return f.inner
}
func (f ListField) Depth() int {
	return f.depth
}

type TypeField struct {
	name  string
	typ   string
	rules []Rule
}

func (f PrimField) Typ() string { return f.typ }
func (f TypeField) Typ() string { return f.typ + "Vgen" }
func (f ListField) Typ() string { return "[]" + f.inner.Typ() }

func (f PrimField) Name() string { return f.name }
func (f TypeField) Name() string { return f.name }
func (f ListField) Name() string { return f.name }

func (f PrimField) LowerName() string { return util.LowerFirstChar(f.name) }
func (f TypeField) LowerName() string { return util.LowerFirstChar(f.name) }
func (f ListField) LowerName() string { return util.LowerFirstChar(f.name) }

// func (f PrimField) HasRules() bool { return len(f.rules.rules) > 0 }
// func (f TypeField) HasRules() bool { return len(f.rules.rules) > 0 }
// func (f ListField) HasRules() bool { return len(f.rules.rules) > 0 }

func (f PrimField) Rules() []Rule { return f.rules }
func (f TypeField) Rules() []Rule { return f.rules }
func (f ListField) Rules() []Rule { return f.rules }

type Tags struct {
	Include bool
}

type Rules struct {
	name     string
	include  bool
	required bool
	rules    [][]Rule
}

// func (r Rules) Name() string      { return r.name }
// func (f Rules) LowerName() string { return util.LowerFirstChar(f.name) }
// func (r Rules) Include() bool     { return r.include }
// func (r Rules) Required() bool    { return r.required }
// func (r Rules) Rules() []Rule {
// 	return r.rules[r.depth]
// }

type Rule struct {
	name  string
	rule  string
	param string
	depth int
}

func (r Rule) Name() string      { return r.name }
func (r Rule) LowerName() string { return util.LowerFirstChar(r.name) }
func (r Rule) Param() string     { return r.param }
func (r Rule) Path() string {
	inner := r.LowerName()
	args := ""
	for i := 0; i < r.depth; i++ {
		inner += "[%d]"
		args += fmt.Sprintf(", i%d", i)
	}
	return fmt.Sprintf(`fmt.Sprintf("%s"%s)`, inner, args)
}

// func (r Rule) Rule() string  { return r.rule }

const DEBUG = false

func parseFile(path string) (ParseInfo, error) {
	// load file
	fset := token.NewFileSet()
	opts := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, path, nil, opts)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not parse file: %v", err)
	}

	// debug print
	if DEBUG {
		err = ast.Print(fset, file)
		if err != nil {
			return ParseInfo{}, fmt.Errorf("could not print ast: %v", err)
		}
	}

	// parse
	var file_err error
	var types []StructType
	ast.Inspect(file, func(n ast.Node) bool {
		if node, ok := n.(*ast.GenDecl); ok {
			parsed_types, err := parseGenDecl(node)
			if err != nil {
				file_err = fmt.Errorf("could not parse gen decl: %v\n", err)
				return false
			}
			for _, s := range parsed_types {
				types = append(types, s)
			}
		}

		return true
	})
	if file_err != nil {
		return ParseInfo{}, file_err
	}

	if DEBUG {
		fmt.Println(types)
	}

	// package
	package_name := file.Name.Name

	return ParseInfo{
		Package:     package_name,
		StructTypes: types,
		Imports:     []string{"encoding/json", "fmt"},
	}, nil

}

func parseGenDecl(node *ast.GenDecl) ([]StructType, error) {
	// check if type
	if node.Tok != token.TYPE {
		return []StructType{}, nil
	}

	if node.Doc != nil {
		tags := parseTags(node.Doc.Text())
		if !tags.Include {
			return []StructType{}, nil
		}
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

	switch t := node.Type.(type) {
	case *ast.StructType:
		return parseStructType(t, node.Name.Name)
	default:
		return StructType{}, fmt.Errorf("unsupported type %T", t)
	}
}

func parseStructType(node *ast.StructType, name string) (StructType, error) {
	if node.Fields == nil || len(node.Fields.List) == 0 {
		return StructType{}, fmt.Errorf("empty structs not supported")
	}

	// parse
	struct_type := StructType{
		name:   name,
		fields: []StructField{},
	}
	for _, field_node := range node.Fields.List {
		field, err := parseStructField(field_node)
		if err != nil {
			return StructType{}, fmt.Errorf("could not parse field: %v", err)
		}
		struct_type.fields = append(struct_type.fields, field)
	}

	return struct_type, nil

}

func parseStructField(node *ast.Field) (StructField, error) {
	name := node.Names[0].Name

	var rules Rules
	var err error
	if node.Comment != nil {
		rules, err = parseRules(node.Comment.Text(), name)
		if err != nil {
			return StructField{}, fmt.Errorf("could not parse rules on %v: %v", name, err)
		}
	}

	// parse
	field, err := parseField(node.Type, name, rules.include, rules.rules, 0)
	if err != nil {
		return StructField{}, fmt.Errorf("could not parse field: %v", err)
	}

	return StructField{
		name:     name,
		field:    field,
		required: rules.required,
	}, nil
}

func parseField(node ast.Expr, name string, include bool, rules [][]Rule, depth int) (Field, error) {
	for len(rules) <= depth {
		rules = append(rules, []Rule{})
	}

	switch n := node.(type) {
	case *ast.Ident:
		typ := n.Name
		if include {
			return TypeField{
				name:  name,
				typ:   typ,
				rules: rules[depth],
			}, nil
		} else {
			return PrimField{
				name:  name,
				typ:   typ,
				rules: rules[depth],
			}, nil
		}
	case *ast.ArrayType:
		inner, err := parseField(n.Elt, name, include, rules, depth+1)
		if err != nil {
			return nil, fmt.Errorf("could not parse arrays inner type: %v", err)
		}
		return ListField{
			name:  name,
			inner: inner,
			rules: rules[depth],
			depth: depth,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported expression: %T", node)
	}
}

func extract(comment string) string {
	reg := regexp.MustCompile(`vgen:\[.+\]`)
	match := reg.FindString(comment)
	match = strings.TrimPrefix(match, "vgen:[")
	match = strings.TrimSuffix(match, "]")
	match = strings.ReplaceAll(match, " ", "")
	return match
}
func parseTags(comment string) Tags {
	extracted := extract(comment)
	var tags Tags
	split := strings.Split(extracted, "|")
	for _, t := range split {
		switch t {
		case "i":
			tags.Include = true
		}
	}
	return tags
}
