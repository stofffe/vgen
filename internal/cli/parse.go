package cli

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"slices"
	"strings"
)

const includeTag = "vgen"

type Parser struct {
	lastPos      token.Pos
	foundImports []string
	fset         *token.FileSet
}

func (p *Parser) SetPos(node ast.Node) {
	p.lastPos = node.Pos()
}

func (p Parser) CurrentLine() *int {
	line := p.fset.Position(p.lastPos).Line
	return &line
}

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

	p := Parser{
		fset:         fset,
		lastPos:      file.Pos(),
		foundImports: []string{},
	}

	// types
	var structTypes []StructType
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// parse tags and check for include
		typeTags, err := p.parseTypeTags(genDecl)
		if err != nil {
			return ParseInfo{}, fmt.Errorf("parse type tags: %w", err)
		}
		if !typeTags.include {
			continue
		}

		// parse types
		for _, spec := range genDecl.Specs {
			spec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, err := p.parseType(spec)
			if err != nil {
				return ParseInfo{}, fmt.Errorf("parse type: %w", err)
			}
			structTypes = append(structTypes, structType)
		}
	}

	// import
	imports := map[string]string{
		"vgen": "github.com/stofffe/vgen/pkg/vgen",
	}
	for _, spec := range file.Imports {
		path := spec.Path.Value
		path = strings.TrimPrefix(path, `"`)
		path = strings.TrimSuffix(path, `"`)
		split := strings.Split(path, "/")
		name := split[len(split)-1]
		// custom name
		if spec.Name != nil {
			name = spec.Name.Name
		}

		if slices.Contains(p.foundImports, name) {
			imports[name] = path
		}
	}

	packageName := file.Name.Name
	return ParseInfo{
		Package:     packageName,
		StructTypes: structTypes,
		Imports:     imports,
	}, nil
}

func (p *Parser) parseType(spec *ast.TypeSpec) (StructType, error) {
	p.SetPos(spec)
	// check name
	if spec.Name == nil {
		return StructType{}, DetailedError{
			line: p.CurrentLine(),
			msg:  "parsed types must have a name",
			err:  fmt.Errorf("must have name"),
		}
	}

	switch node := spec.Type.(type) {
	case *ast.StructType:
		name := spec.Name.Name
		structType, err := p.parseStruct(node, name)
		if err != nil {
			return StructType{}, fmt.Errorf("prase struct: %w", err)
		}
		return structType, nil
	case *ast.Ident:
		return StructType{}, DetailedError{
			line: p.CurrentLine(),
			msg:  "type aliases not supported",
			err:  fmt.Errorf("type aliases not supported"),
		}
	default:
		return StructType{}, fmt.Errorf("unsupported type %T", node)
	}
}

func (p *Parser) parseStruct(structNode *ast.StructType, structName string) (StructType, error) {
	p.SetPos(structNode)
	structType := StructType{
		Name:   structName,
		Fields: []StructField{},
	}

	for _, fieldNode := range structNode.Fields.List {
		field, err := p.parseField(fieldNode)
		if err != nil {
			return StructType{}, fmt.Errorf("parse field: %w", err)
		}
		structType.Fields = append(structType.Fields, field)
	}

	return structType, nil
}

func (p *Parser) parseField(fieldNode *ast.Field) (StructField, error) {
	p.SetPos(fieldNode)
	comments := fieldNode.Doc.Text() + fieldNode.Comment.Text()

	fieldName := fieldNode.Names[0].Name // TODO handle multiple
	tags := ""
	if fieldNode.Tag != nil {
		tags = fieldNode.Tag.Value
	}

	// parse tags
	commentTags, err := p.parseFieldTags(comments)
	if err != nil {
		return StructField{}, fmt.Errorf("parse field tags: %w", err)
	}
	nested := commentTags.include
	alias := fieldName
	if commentTags.alias != "" {
		alias = commentTags.alias
	}

	// parse field types
	fieldInfo, err := p.parseFieldType(fieldNode.Type)
	if err != nil {
		return StructField{}, fmt.Errorf("field type: %w", err)
	}

	// dont allow nested on primitve types
	if nested && fieldInfo.Primitive {
		return StructField{}, DetailedError{
			line: p.CurrentLine(),
			msg:  "primitve fields can not have nested tag",
			err:  fmt.Errorf("nested not allowed on primitve inner type"),
		}
	}

	return StructField{
		Name:   fieldName,
		Tags:   tags,
		Alias:  alias,
		Nested: nested,

		Types:   fieldInfo.Types,
		Pointer: fieldInfo.Pointer,
		Import:  fieldInfo.Import,
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

func (p *Parser) parseFieldType(node ast.Expr) (FieldTypeInfo, error) {
	info := FieldTypeInfo{
		Types:     []FieldType{},
		Pointer:   false,
		Primitive: false,
		Import:    "",
	}
	currentNode := node
	for {
		switch node := currentNode.(type) {
		// internal primitive/struct
		case *ast.Ident:
			if node.Obj == nil {
				info.Primitive = true
			}
			info.Types = append(info.Types, FieldTypeIdent{name: node.Name})
			return info, nil
		// import
		case *ast.SelectorExpr:
			imp, ok := node.X.(*ast.Ident)
			if !ok {
				return FieldTypeInfo{}, fmt.Errorf("import selector is not ast.Ident")
			}
			p.foundImports = append(p.foundImports, imp.Name)
			info.Import = imp.Name
			info.Types = append(info.Types, FieldTypeImport{imp: imp.Name, name: node.Sel.Name})
			return info, nil
		// array
		case *ast.ArrayType:
			info.Types = append(info.Types, FieldTypeArray{})
			currentNode = node.Elt
		// map
		case *ast.MapType:
			info.Types = append(info.Types, FieldTypeMap{})
			key, ok := node.Key.(*ast.Ident)
			if !ok || key.Name != "string" {
				return FieldTypeInfo{}, DetailedError{
					line: p.CurrentLine(),
					msg:  "key of map must be a string",
					err:  fmt.Errorf("invalid map key %v, must be string", key),
				}
			}
			currentNode = node.Value
		// pointer
		case *ast.StarExpr:
			if len(info.Types) > 0 {
				return FieldTypeInfo{}, DetailedError{
					line: p.CurrentLine(),
					msg:  "pointers not allowed as list/map element",
					err:  fmt.Errorf("pointers not allowed as list/map element"),
				}
			}
			info.Pointer = true
			currentNode = node.X
		default:
			return FieldTypeInfo{}, fmt.Errorf("unknown field type %T", node)
		}
	}
}

type FieldTags struct {
	include bool
	alias   string
}

func (p *Parser) parseFieldTags(comment string) (FieldTags, error) {
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

	if len(args) == 1 && args[0] == "" {
		return FieldTags{}, DetailedError{
			line: p.CurrentLine(),
			msg:  "field tags can not be empty",
			err:  fmt.Errorf("field tag empty"),
		}
	}

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
					line: p.CurrentLine(),
					msg:  `field tag alias must have a value, ex: "alias=something"`,
					err:  fmt.Errorf("alias must have second argument"),
				}
			}
			name := split[1]
			tags.alias = name
		default:
			return FieldTags{}, DetailedError{
				line: p.CurrentLine(),
				msg:  fmt.Sprintf("unknown field tag %v", ident),
				err:  fmt.Errorf("unknown field tag %v", ident),
			}
		}
	}
	return tags, nil
}

type TypeTags struct {
	include bool
}

func (p *Parser) parseTypeTags(genDecl *ast.GenDecl) (TypeTags, error) {
	p.SetPos(genDecl)

	comment := genDecl.Doc.Text()
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

	if len(args) == 1 && args[0] == "" {
		return TypeTags{}, DetailedError{
			line: p.CurrentLine(),
			msg:  "type tags can not be empty",
			err:  fmt.Errorf("type tag empty"),
		}
	}

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		split := strings.Split(arg, "=")
		ident := split[0]
		switch ident {
		case "include", "i":
			tags.include = true
		default:
			return TypeTags{}, DetailedError{
				line: p.CurrentLine(),
				msg:  fmt.Sprintf("unknown type tag %s", ident),
				err:  fmt.Errorf("unknown type tag %s", ident),
			}
		}
	}
	return tags, nil
}
