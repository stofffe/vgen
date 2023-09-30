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
	Structs []Type
}

// types
type Type interface {
	isType()
}
type InvalidType struct{}

func (t InvalidType) isType() {}

type StructType struct {
	Name   string
	Fields []Field
}

func (t StructType) isType() {}

// field
type Field interface {
	isField()
}
type PrimitiveField struct {
	Name     string
	Typ      string
	Rules    []Rule
	Required bool
}

func (f PrimitiveField) isField() {}

type ListField struct {
	Name       string
	Typ        string
	ListRules  []Rule
	ValueRules []Rule
	Required   bool
}

func (f ListField) isField() {}

type InvalidField struct{}

func (f InvalidField) isField() {}

// rule
type Rule interface {
	isRule()
}

type PrimitiveRule struct {
	FieldName    string
	Func         string
	Value        string
	IncludeIndex bool
}

func (t PrimitiveRule) isRule() {}

func parseFile(path string) (ParseInfo, error) {
	// load file
	fset := token.NewFileSet()
	opts := parser.AllErrors | parser.ParseComments
	file, err := parser.ParseFile(fset, path, nil, opts)
	if err != nil {
		return ParseInfo{}, fmt.Errorf("could not parse file: %v", err)
	}

	// err = ast.Print(fset, file)
	// if err != nil {
	// 	return ParseInfo{}, fmt.Errorf("could not print ast: %v", err)
	// }

	// parse
	var file_err error
	var structs []Type
	ast.Inspect(file, func(n ast.Node) bool {
		if node, ok := n.(*ast.GenDecl); ok {
			parsed_structs, err := parseGenDecl(node)
			if err != nil {
				file_err = fmt.Errorf("could not parse gen decl: %v\n", err)
				return false
			}
			for _, s := range parsed_structs {
				// fmt.Println(s)
				structs = append(structs, s)
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
		Structs: structs,
	}, nil

}

const INCLUDE_TAG = `i`

func parseGenDecl(node *ast.GenDecl) ([]Type, error) {
	// check for tag
	if node.Doc == nil {
		return []Type{}, nil
	}
	hasTag := false
	for _, comment := range node.Doc.List {
		rules := extractRules(comment.Text)
		for _, rule := range rules {
			if rule == "i" {
				hasTag = true
			}
		}
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
		return InvalidType{}, fmt.Errorf("must have name")
	}

	switch t := node.Type.(type) {
	case *ast.StructType:
		return parseStructType(t, node.Name.Name)
	default:
		return InvalidType{}, fmt.Errorf("unsupported type %T", t)
	}
}

func parseStructType(node *ast.StructType, name string) (Type, error) {
	if node.Fields == nil || len(node.Fields.List) == 0 {
		return InvalidType{}, fmt.Errorf("empty structs not supported")
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
		return InvalidField{}, fmt.Errorf("field without name not supported")
	}
	if len(node.Names) > 1 {
		return InvalidField{}, fmt.Errorf("field with multiple names not supported")
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
	var err error

	switch n := node.Type.(type) {
	case *ast.Ident:
		field, err = parsePrimitiveField(n, name, comment)
		if err != nil {
			return InvalidField{}, fmt.Errorf("could not parse field primtive: %v", err)
		}
	case *ast.ArrayType:
		field, err = parseListField(n, name, comment)
		if err != nil {
			return InvalidField{}, fmt.Errorf("could not parse field primtive: %v", err)
		}
	default:
		return InvalidField{}, fmt.Errorf("unsupported field type: %T", n)
	}

	return field, nil
}

// no nested lists
func parseListField(node *ast.ArrayType, field_name, comment string) (Field, error) {
	// type
	inner_type, ok := node.Elt.(*ast.Ident)
	if !ok {
		return InvalidField{}, fmt.Errorf("type of array must be primitive not: %T", node.Elt)
	}
	typ := "[]" + inner_type.Name

	// rules
	rules_str := extractRules(comment)
	var list_rules_str []string
	var value_rules_str []string
	req := false
	for _, rule := range rules_str {
		if rule == "req" {
			req = true
			continue
		}

		if rule == "" {
			return InvalidField{}, fmt.Errorf("empty rule")
		}

		// check rule
		if []rune(rule)[0] == ':' {
			value_rules_str = append(value_rules_str, rule[1:])
		} else {
			list_rules_str = append(list_rules_str, rule)
		}
	}

	list_rules, err := parseRules(list_rules_str, field_name, false)
	if err != nil {
		return InvalidField{}, fmt.Errorf("could not parse list rules: %v", err)
	}
	value_rules, err := parseRules(value_rules_str, field_name, true)
	if err != nil {
		return InvalidField{}, fmt.Errorf("could not parse value rules: %v", err)
	}

	return ListField{
		Name:       field_name,
		Typ:        typ,
		ListRules:  list_rules,
		ValueRules: value_rules,
		Required:   req,
	}, nil
}

func parsePrimitiveField(node *ast.Ident, field_name, comment string) (Field, error) {
	// rules
	rules_str := extractRules(comment)
	req := false
	for _, rule := range rules_str {
		if rule == "req" {
			req = true
			break
		}
	}
	rules, err := parseRules(rules_str, field_name, false)
	if err != nil {
		return InvalidField{}, fmt.Errorf("could not parse rules: %v", err)
	}

	typ := node.Name
	return PrimitiveField{
		Name:     field_name,
		Typ:      typ,
		Rules:    rules,
		Required: req,
	}, nil
}

var extractRulesRegex = createExtractRulesRegex()

func createExtractRulesRegex() *regexp.Regexp {
	return regexp.MustCompile(`vgen:\[(.*)\]`)
}

func extractRules(value string) []string {
	matches := extractRulesRegex.FindStringSubmatch(value)

	if len(matches) == 0 {
		return []string{}
	}

	rules := matches[1] // first match is whole string
	rules = strings.ReplaceAll(rules, " ", "")
	split := strings.Split(rules, ",")

	return split
}

var parseRulesRegex = createParseRulesRegex()

func createParseRulesRegex() *regexp.Regexp {
	req := `^(req)$`                 // all
	len_gt := `^(len_gt)\((.+)\)$`   // string, list, map
	len_lt := `^(len_lt)\((.+)\)$`   // string, list, map
	len_gte := `^(len_gte)\((.+)\)$` // string, list, map
	len_lte := `^(len_lte)\((.+)\)$` // string, list, map
	not_empty := `^(not_empty)$`     // string, list, map
	gt := `^(gt)\((.+)\)$`           // string, int, float
	lt := `^(lt)\((.+)\)$`           // string, int, float
	gte := `^(gte)\((.+)\)$`         // string, int, float
	lte := `^(lte)\((.+)\)$`         // string, int, float
	custom := `^(custom)\((.+)\)$`   // all

	rules := []string{req, len_gt, len_lt, len_gte, len_lte, not_empty, gt, lt, gte, lte, custom}
	pattern := strings.Join(rules, "|")

	return regexp.MustCompile(pattern)
}

func parseRules(rules_str []string, name string, include_index bool) ([]Rule, error) {
	var rules []Rule

	for _, rule := range rules_str {
		matches := parseRulesRegex.FindStringSubmatch(rule)

		if len(matches) == 0 {
			return []Rule{}, fmt.Errorf("invalid rule: %v", rule)
		}

		var filtered []string
		for i, v := range matches {
			if i != 0 && v != "" {
				filtered = append(filtered, v)
			}
		}

		// extract func and parameter (if exists)
		f := filtered[0]
		v := ""
		if len(filtered) > 1 {
			v = filtered[1]
		}

		// TODO add custom fieldname if json:"" supplied?
		rules = append(rules, PrimitiveRule{
			FieldName:    name,
			Func:         f,
			Value:        v,
			IncludeIndex: include_index,
		})
	}

	return rules, nil

}
