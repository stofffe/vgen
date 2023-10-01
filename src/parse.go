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
	Types   []Type
	Imports []string
}

// types
type Type interface {
	isType()
}
type StructType struct {
	Name   string
	Fields []Field
}

func (t StructType) isType() {}

// field
type Field interface {
	isField()
	Typ() string
	FieldValidationCode() (string, error)
}
type PrimitiveField struct {
	Name     string
	typ      string
	Rules    []Rule
	Required bool
}
type ListField struct {
	Name       string
	innerType  string
	ListRules  []Rule
	ValueRules []Rule
	Required   bool
}
type TypeField struct {
	Name     string
	typ      string
	Rules    []Rule
	Required bool
}

func (f PrimitiveField) Typ() string {
	return f.typ
}
func (f TypeField) Typ() string {
	return f.typ + "Vgen"
}
func (f ListField) Typ() string {
	return "[]" + f.innerType
}

func (f PrimitiveField) isField() {}
func (f ListField) isField()      {}
func (f TypeField) isField()      {}

// rule
type Rule interface {
	isRule()
	RuleValidationCode() (string, error)
}
type PrimitiveRule struct {
	FieldName string
	Func      string
	Value     string
}
type ListRule struct {
	FieldName string
	Func      string
	Value     string
}

func (t ListRule) isRule()      {}
func (t PrimitiveRule) isRule() {}

type Tags struct {
	Include bool
}

type Rules struct {
	FieldName      string
	Typ            string
	Include        bool
	Required       bool
	PrimitiveRules []Rule
	ListRules      []Rule
}

func parseFile(path string) (ParseInfo, error) {
	// load file
	fset := token.NewFileSet()
	opts := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, path, nil, opts)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not parse file: %v", err)
	}

	// debug print
	err = ast.Print(fset, file)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not print ast: %v", err)
	}

	// parse
	var file_err error
	var types []Type
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

	// package
	package_name := file.Name.Name

	return ParseInfo{
		Package: package_name,
		Types:   types,
		Imports: []string{"encoding/json", "fmt"},
	}, nil

}

func parseGenDecl(node *ast.GenDecl) ([]Type, error) {
	// check for tag
	if node.Doc == nil {
		return []Type{}, nil
	}
	hasTag := false
	for _, comment := range node.Doc.List {
		// TODO split func into tags and rules?
		tags, err := parseTags(comment.Text)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for type: %v", tags)
		}
		hasTag = tags.Include
	}
	if !hasTag {
		return []Type{}, nil
	}

	// check if type
	if node.Tok != token.TYPE {
		return []Type{}, nil
	}

	// parse all types under decl
	var structs []Type
	for _, spec := range node.Specs {
		type_node := spec.(*ast.TypeSpec) // already checked

		s, err := parseType(type_node)
		structs = append(structs, s)
		if err != nil {
			return []Type{}, fmt.Errorf("could not parse type: %v", err)
		}
	}

	return structs, nil
}

func parseType(node *ast.TypeSpec) (Type, error) {
	// check name
	if node.Name == nil {
		return nil, fmt.Errorf("must have name")
	}

	switch t := node.Type.(type) {
	case *ast.StructType:
		return parseStructType(t, node.Name.Name)
	default:
		return nil, fmt.Errorf("unsupported type %T", t)
	}
}

func parseStructType(node *ast.StructType, name string) (Type, error) {
	if node.Fields == nil || len(node.Fields.List) == 0 {
		return nil, fmt.Errorf("empty structs not supported")
	}

	// parse
	struct_type := StructType{
		Name:   name,
		Fields: []Field{},
	}
	for _, field_node := range node.Fields.List {
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
		return nil, fmt.Errorf("field without name not supported")
	}
	if len(node.Names) > 1 {
		return nil, fmt.Errorf("field with multiple names not supported")
	}
	name := node.Names[0].Name

	// extract rules from comment
	var comment string
	if node.Comment != nil {
		if len(node.Comment.List) > 1 {
			log.Fatalf("HOW CAN THIS BE > 1?")
		}
		comment = node.Comment.List[0].Text
	}

	// parse
	var field Field
	switch n := node.Type.(type) {
	case *ast.Ident:
		typ := n.Name
		rules, err := parseRules(name, typ, comment)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for field: %v", err)
		}

		if rules.Include {
			field = TypeField{
				Name:     name,
				typ:      n.Name,
				Rules:    rules.PrimitiveRules,
				Required: rules.Required,
			}
		} else {
			field = PrimitiveField{
				Name:     name,
				typ:      typ,
				Rules:    rules.PrimitiveRules,
				Required: rules.Required,
			}
		}

	case *ast.ArrayType:
		// type
		inner_type, ok := n.Elt.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("type of array must be primitive not: %T", n.Elt)
		}
		rules, err := parseRules(name, inner_type.Name, comment)
		if err != nil {
			return nil, fmt.Errorf("invalid tags for field: %v", err)
		}
		field = ListField{
			Name:       rules.FieldName,
			innerType:  inner_type.Name,
			ListRules:  rules.PrimitiveRules,
			ValueRules: rules.ListRules,
			Required:   rules.Required,
		}

	default:
		return nil, fmt.Errorf("unsupported field type: %T", n)
	}

	return field, nil
}

func parseRules(name, typ, comment string) (Rules, error) {
	// extract tags
	reg := regexp.MustCompile(`vgen:\[(.+)\]`)
	matches := reg.FindStringSubmatch(comment)

	if len(matches) == 0 {
		return Rules{}, nil
	}

	content := matches[1]
	content = strings.ReplaceAll(content, " ", "") // first match is whole string
	split := strings.Split(content, ",")

	tags := Rules{
		FieldName: name,
		Typ:       typ,
	}
	for _, tag := range split {
		split := strings.Split(tag, "=")
		rule := strings.ReplaceAll(split[0], " ", "")
		param := ""
		if len(split) > 1 {
			param = strings.ReplaceAll(split[1], " ", "")

		}

		switch rule {
		// special rule
		case "i":
			tags.Include = true
			continue
		// special rule
		case "req":
			tags.Required = true
			continue
		case "not_empty", "custom", "len_gt", "len_gte", "len_lt", "len_lte", "gt", "gte", "lt", "lte":
			// if !(typ == "string" || strings.HasPrefix(rule, "[]")) {
			// 	return Tags{}, fmt.Errorf("")
			// }
			tags.PrimitiveRules = append(tags.PrimitiveRules, PrimitiveRule{
				FieldName: name,
				Func:      rule,
				Value:     param,
			})
			continue
		case ":not_empty", ":custom", ":len_gt", ":len_gte", ":len_lt", ":len_lte", ":gt", ":gte", ":lt", ":lte":
			tags.ListRules = append(tags.ListRules, ListRule{
				FieldName: name,
				Func:      rule,
				Value:     param,
			})
			continue
		default:
			return Rules{}, fmt.Errorf("invalid rule: %v", rule)
		}

	}
	return tags, nil
}

func parseTags(comment string) (Tags, error) {
	split := extract(comment)

	var tags Tags
	for _, tag := range split {
		switch tag {
		// special rule
		case "i":
			tags.Include = true
			continue
		default:
			return Tags{}, fmt.Errorf("invalid tag: %v", tag)
		}
	}
	return tags, nil
}

func extract(comment string) []string {
	// extract tags
	reg := regexp.MustCompile(`vgen:\[(.+)\]`)
	matches := reg.FindStringSubmatch(comment)
	if len(matches) == 0 {
		return []string{}
	}

	content := matches[1]
	content = strings.ReplaceAll(content, " ", "") // first match is whole string
	return strings.Split(content, ",")

}
