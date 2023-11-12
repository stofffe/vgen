package main

import (
	"fmt"
	"regexp"
	"unicode"
)

// String to rules

func parseRules(input, name, alias string, pos int) (Rules, error) {
	// remove comment
	input = input[2:]
	pos += 2

	lexer := NexLexer(input, pos)
	lexer.Lex()

	parser := NewParser(lexer, name, alias)
	err := parser.Parse()
	if err != nil {
		return Rules{}, fmt.Errorf("could not parse rules: %v", err)
	}

	return parser.rules, nil
}

//
// Lexer
//

type Lexer struct {
	file_pos int
	input    string
	tokens   []Token
}

func NexLexer(input string, pos int) Lexer {
	// regex match neccesary info
	reg := regexp.MustCompile(`(.*vgen:)(\[.+\])(.*)`)
	matches := reg.FindStringSubmatch(input)
	if len(matches) == 0 {
		return Lexer{}
	}

	prefix := matches[1]
	match := matches[2]

	pos += len(prefix)

	return Lexer{
		file_pos: pos,
		input:    match,
		tokens:   []Token{},
	}
}

func (l *Lexer) AddToken(token Token) {
	l.tokens = append(l.tokens, token)
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	EQUAL
	IDENT
	NUMBER
	COLON
	QUOTE
	SPACE
)

var TOKEN_NAMES = map[TokenType]string{
	ILLEGAL:     "ILLEGAL",
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COMMA:       "COMMA",
	EQUAL:       "EQUAL",
	IDENT:       "IDENT",
	NUMBER:      "NUMBER",
	COLON:       "COLON",
	QUOTE:       "QUOTE",
	SPACE:       "SPACE",
}

func (t TokenType) String() string {
	return TOKEN_NAMES[t]
}

type Token struct {
	typ   TokenType
	value string
	pos   int
}

func (t Token) FilePosition() {

}

func (l *Lexer) Lex() {
	for len(l.input) > 0 {
		// Number TODO negative numbers
		if unicode.IsDigit(l.Peek()) {
			l.LexNumber()
			continue
		}

		// Rules
		if unicode.IsLetter(l.Peek()) {
			l.LexString()
			continue
		}

		// Custom error messages
		// TODO
		if l.Peek() == '"' {

		}

		// Single chars
		n := l.Consume()
		switch n {
		case '[':
			l.AddToken(Token{typ: LEFT_BRACE, pos: l.file_pos})
		case ']':
			l.AddToken(Token{typ: RIGHT_BRACE, pos: l.file_pos})
		case ',':
			l.AddToken(Token{typ: COMMA, pos: l.file_pos})
		case '=':
			l.AddToken(Token{typ: EQUAL, pos: l.file_pos})
		case ':':
			l.AddToken(Token{typ: COLON, pos: l.file_pos})
		case ' ', '\t':
			l.AddToken(Token{typ: SPACE, pos: l.file_pos})
		default:
			l.AddToken(Token{typ: ILLEGAL, value: string(n), pos: l.file_pos})
			return
		}
		l.file_pos++
	}
}

func (l *Lexer) LexQuote() {

}

func (l *Lexer) LexNumber() {
	pos := l.file_pos
	number := []rune{}
	for l.InputLeft() && unicode.IsDigit(l.Peek()) {
		number = append(number, l.Consume())
		l.file_pos++
	}
	l.AddToken(Token{
		typ:   NUMBER,
		value: string(number),
		pos:   pos,
	})
}

func (l *Lexer) LexString() {
	pos := l.file_pos
	rule := []rune{}
	for l.InputLeft() && (l.Peek() == '_' || unicode.IsLetter(l.Peek())) {
		rule = append(rule, l.Consume())
		l.file_pos++
	}
	l.AddToken(Token{
		typ:   IDENT,
		value: string(rule),
		pos:   pos,
	})
}

func (l *Lexer) InputLeft() bool {
	return len(l.input) > 0
}

// TODO no range check
func (l *Lexer) Peek() rune {
	return rune(l.input[0])
}

// TODO no range check
func (l *Lexer) Consume() rune {
	r := rune(l.input[0])
	l.input = l.input[1:]
	return r
}

func PrintTokens(tokens []Token) {
	for _, token := range tokens {
		fmt.Printf("%s\t%s\t%d\n", TOKEN_NAMES[token.typ], token.value, token.pos)
	}
}

//
// Parser
//

type Parser struct {
	tokens []Token
	rules  Rules
	depth  int
	name   string
	alias  string
}

func NewParser(lexer Lexer, name string, alias string) Parser {
	return Parser{
		name:   name,
		alias:  alias,
		tokens: lexer.tokens, rules: Rules{
			name:     name,
			include:  false,
			required: false,
			rules:    [][]Rule{},
		}, depth: -1}
}

func (p *Parser) AddRule(rule Rule) {
	for len(p.rules.rules) <= p.depth {
		p.rules.rules = append(p.rules.rules, []Rule{})
	}
	p.rules.rules[p.depth] = append(p.rules.rules[p.depth], rule)
}

func (p *Parser) TokensLeft() bool {
	return len(p.tokens) > 0
}

// TODO add ok check to this
func (p *Parser) Peek() Token {
	return p.tokens[0]
}

func (p *Parser) PeekTyp(typ TokenType) bool {
	return p.Peek().typ == typ
}

// TODO add ok check to this
func (p *Parser) Consume() Token {
	r := p.tokens[0]
	p.tokens = p.tokens[1:]
	return r
}

// Consumes zero or more spaces
func (p *Parser) consumeSpaces() {
	for p.TokensLeft() && p.PeekTyp(SPACE) {
		p.Consume()
	}
}

func (p *Parser) Parse() error {
	return p.expectRules()
}

func (p *Parser) expectRules() error {
	// increase depth
	p.depth++

	// empty
	if !p.TokensLeft() {
		return nil
	}

	// [<rule>]<rules>
	_, err := p.expectToken(LEFT_BRACE)
	if err != nil {
		return err
	}

	err = p.expectRule()
	if err != nil {
		return err
	}

	_, err = p.expectToken(RIGHT_BRACE)
	if err != nil {
		return err
	}

	return p.expectRules()
}

func (p *Parser) expectRule() error {
	p.consumeSpaces()

	// handle empty
	if p.PeekTyp(RIGHT_BRACE) {
		return nil
	}

	token, err := p.expectToken(IDENT)
	if err != nil {
		return err
	}

	rule := token.value
	switch rule {
	// Special rule
	case "req", "required", "i", "include":
		err := p.expectNoArgRule(rule)
		if err != nil {
			return err
		}
	// No args rule
	case "not_empty":
		err := p.expectNoArgRule(rule)
		if err != nil {
			return err
		}
	// Decimal rule
	case "gt", "lt", "gte", "lte", "len_gt", "len_gte", "len_lt", "len_lte":
		err := p.expectDecimalRule(rule)
		if err != nil {
			return err
		}
	// Ident rule
	case "custom":
		err := p.expectIdentRule(rule)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected rule %s", rule)
	}

	if p.PeekTyp(COMMA) {
		p.Consume()
		return p.expectRule()
	}

	return nil
}

func (p *Parser) expectNoArgRule(rule string) error {
	switch rule {
	case "req", "required":
		p.rules.required = true
	case "i", "include":
		p.rules.include = true
		p.AddRule(Rule{
			name:  p.name,
			alias: p.alias,
			rule:  rule,
			depth: p.depth,
		})
	default:
		p.AddRule(Rule{
			name:  p.name,
			alias: p.alias,
			rule:  rule,
			depth: p.depth,
		})
	}

	p.consumeSpaces()
	return nil
}

func (p *Parser) expectDecimalRule(rule string) error {
	_, err := p.expectToken(EQUAL)
	if err != nil {
		return err
	}
	token, err := p.expectToken(NUMBER)
	if err != nil {
		return err
	}
	p.AddRule(Rule{
		name:  p.name,
		alias: p.alias,
		rule:  rule,
		param: token.value,
		depth: p.depth,
	})

	p.consumeSpaces()
	return nil
}

func (p *Parser) expectIdentRule(rule string) error {
	_, err := p.expectToken(EQUAL)
	if err != nil {
		return err
	}
	token, err := p.expectToken(IDENT)
	if err != nil {
		return err
	}
	p.AddRule(Rule{
		name:  p.name,
		alias: p.alias,
		rule:  rule,
		param: token.value,
		depth: p.depth,
	})

	p.consumeSpaces()
	return nil
}

func (p *Parser) expectToken(expected TokenType) (Token, error) {
	token := p.Consume()
	if token.typ != expected {
		return Token{}, fmt.Errorf("unexpected token at byte pos %d: expected %s got %s", token.pos, expected.String(), token.typ.String())
	}
	return token, nil
}
