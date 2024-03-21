package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

const DEBUG = false

type ParseInfo struct {
	Package     string
	StructTypes []StructType
	Imports     []string
}

type StructType struct {
	name   string
	fields []StructField
}

func (s StructType) Name() string {
	return s.name
}

func (s StructType) Fields() []StructField {
	return s.fields
}

type StructField struct {
	name     string
	alias    string
	include  bool
	required bool
	innerTyp string
	depth    int
	tags     string
	rules    []Rule
	// field    Field
}

func (s StructField) Name() string     { return s.name }
func (s StructField) Alias() string    { return s.alias }
func (s StructField) Include() bool    { return s.include }
func (s StructField) Required() bool   { return s.required }
func (s StructField) InnerTyp() string { return s.innerTyp }
func (s StructField) Tags() string     { return s.tags }
func (s StructField) Depth() int       { return s.depth }
func (s StructField) Rules() []Rule    { return s.rules }
func (s StructField) Typ() string {
	var buffer strings.Builder
	for i := 0; i < s.depth; i++ {
		buffer.WriteString("[]")
	}
	buffer.WriteString(s.innerTyp)
	if s.include {
		buffer.WriteString("Vgen")
	}
	return buffer.String()
}
func (s StructField) ConvTyp() string {
	var buffer strings.Builder
	for i := 0; i < s.depth; i++ {
		buffer.WriteString("[]")
	}
	buffer.WriteString(s.innerTyp)
	return buffer.String()
}

type Tags struct {
	Include bool
}

type Rules struct {
	include  bool
	required bool
	rules    []Rule
}

type Rule struct {
	rule  string
	param string
	field *StructField
}

func (r Rule) Rule() string        { return r.rule }
func (r Rule) Param() string       { return r.param }
func (r Rule) Field() *StructField { return r.field }

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
		// Imports:     []string{"encoding/json", "fmt"},
		Imports: []string{"fmt"},
	}, nil

}

func parseGenDecl(node *ast.GenDecl) ([]StructType, error) {
	// check if type
	if node.Tok != token.TYPE {
		return []StructType{}, nil
	}

	// Parse type tags
	var tags Tags
	if node.Doc != nil {
		tags = parseTags(node.Doc.Text())
	}
	if node.Doc == nil || !tags.Include {
		return []StructType{}, nil
	}

	// parse all types under decl
	var structs []StructType
	for _, spec := range node.Specs {
		type_node := spec.(*ast.TypeSpec) // already checked

		s, err := parseType(type_node)
		if err != nil {
			return []StructType{}, fmt.Errorf("could not parse type: %v", err)
		}
		structs = append(structs, s)
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
	// empty check
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
	// field tags, ex json:""
	tags := ""
	if node.Tag != nil {
		tags = node.Tag.Value
	}

	// name and alias
	name := node.Names[0].Name
	alias := ExtractJsonName(tags, name)

	// rules
	var rules Rules
	var err error
	if node.Comment != nil {
		rules, err = parseRules(extract(node.Comment.Text()))
		if err != nil {
			return StructField{}, fmt.Errorf("could not parse rules: %v", err)
		}
		// fmt.Println(alias, s)
	}

	if len(rules.rules) == 0 && !rules.include && !rules.required {
		return StructField{}, fmt.Errorf("empty fields not allowed")
	}

	// parse
	depth, inner, err := getDepth(node.Type, 0)
	if err != nil {
		return StructField{}, fmt.Errorf("could not get depth: %v", err)
	}

	result := StructField{
		name:     name,
		alias:    alias,
		include:  rules.include,
		required: rules.required,
		depth:    depth,
		tags:     tags,
		rules:    rules.rules,
		innerTyp: inner,
	}
	// parent pointers
	for i := range result.rules {
		result.rules[i].field = &result
	}

	return result, nil
}

func getDepth(node ast.Expr, depth int) (int, string, error) {
	switch n := node.(type) {
	case *ast.Ident:
		return depth, n.Name, nil
	case *ast.ArrayType:
		depth, inner, err := getDepth(n.Elt, depth+1)
		if err != nil {
			return 0, "", fmt.Errorf("could not parse arrays inner type: %v", err)
		}
		return depth, inner, nil
	default:
		return 0, "", fmt.Errorf("unsupported expression: %T", node)
	}
}

func extract(comment string) string {
	reg := regexp.MustCompile(`vgen:\[[^\]]+\]`)
	match := reg.FindString(comment)
	match = strings.TrimPrefix(match, "vgen:[")
	match = strings.TrimSuffix(match, "]")
	match = strings.ReplaceAll(match, " ", "")
	return match
}

func parseRules(comment string) (Rules, error) {
	rules := Rules{
		include:  false,
		required: false,
		rules:    []Rule{},
	}

	for _, rule := range strings.Split(comment, ",") {
		// No args
		switch rule {
		// Special rule
		case "req", "required":
			rules.required = true
			continue
		case "i", "include":
			rules.include = true
			continue
		// NoArg rule
		case "not_empty":
			rules.rules = append(rules.rules, Rule{
				rule: rule,
			})
			continue
		}
		// Args
		parts := strings.Split(rule, "=")
		if len(parts) == 1 {
			return Rules{}, fmt.Errorf("invalid rule \"%v\"", rule)
		}
		switch parts[0] {
		case "gt", "lt", "gte", "lte", "len_gt", "len_gte", "len_lt", "len_lte", "custom":
			// TODO validate number
			rules.rules = append(rules.rules, Rule{
				rule:  parts[0],
				param: parts[1],
			})
		default:
			return Rules{}, fmt.Errorf("unexpected rule %s", rule)
		}
	}

	return rules, nil
}

func parseTags(comment string) Tags {
	extracted := extract(comment)
	var tags Tags
	for _, tag := range strings.Split(extracted, ",") {
		switch tag {
		case "i":
			tags.Include = true
		}
	}
	return tags
}

func ExtractJsonName(tag, backup string) string {
	reg := regexp.MustCompile(`json:"[^"]*"`)
	match := reg.FindString(tag)
	match = strings.TrimPrefix(match, `json:"`)
	match = strings.TrimSuffix(match, `"`)
	match = strings.Split(match, ",")[0]
	if match == "" {
		return backup
	}
	return match
}
