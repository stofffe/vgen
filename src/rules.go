package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// String to rules

func parseRules(input string, name string) (Rules, error) {
	lexer := NexLexer(input)
	lexer.Lex()

	parser := NewParser(lexer, name)
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
	pos    int
	input  string
	tokens []Token
}

func NexLexer(input string) Lexer {
	// remove whitespace
	input = strings.ReplaceAll(input, " ", "")

	// regex match neccesary info
	reg := regexp.MustCompile(`vgen:(\[.+\])`)
	matches := reg.FindStringSubmatch(input)
	if len(matches) == 0 {
		return Lexer{}
	}
	match := matches[1]

	return Lexer{
		pos:    0,
		input:  match,
		tokens: []Token{},
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
)

var TOKEN_NAMES = map[TokenType]string{
	ILLEGAL:     "ILLEGAL",
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COMMA:       "COMMA",
	EQUAL:       "EQUAL",
	IDENT:       "IDENT",
	NUMBER:      "NUMBER",
}

func (t TokenType) String() string {
	return TOKEN_NAMES[t]
}

type Token struct {
	typ   TokenType
	value string
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
			l.LexRule()
			continue
		}

		// Single chars
		n := l.Consume()
		switch n {
		case '[':
			l.AddToken(Token{typ: LEFT_BRACE})
		case ']':
			l.AddToken(Token{typ: RIGHT_BRACE})
		case ',':
			l.AddToken(Token{typ: COMMA})
		case '=':
			l.AddToken(Token{typ: EQUAL})
		default:
			l.AddToken(Token{typ: ILLEGAL, value: string(n)})
			return
		}
	}
}

func (l *Lexer) LexNumber() {
	number := []rune{}
	for l.InputLeft() && unicode.IsDigit(l.Peek()) {
		number = append(number, l.Consume())
	}
	l.AddToken(Token{
		typ:   NUMBER,
		value: string(number),
	})
}

func (l *Lexer) LexRule() {
	rule := []rune{}
	for l.InputLeft() && (l.Peek() == '_' || unicode.IsLetter(l.Peek())) {
		rule = append(rule, l.Consume())
	}
	l.AddToken(Token{
		typ:   IDENT,
		value: string(rule),
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
		fmt.Printf("%s\t%s\n", TOKEN_NAMES[token.typ], token.value)
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
}

func NewParser(lexer Lexer, name string) Parser {
	return Parser{
		name:   name,
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

// TODO add ok check to this
func (p *Parser) Consume() Token {
	r := p.tokens[0]
	p.tokens = p.tokens[1:]
	return r
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
	token, err := p.expectToken(IDENT)
	if err != nil {
		return err
	}

	rule := token.value
	switch rule {
	case "req", "required", "i", "include":
		err := p.expectNoArgRule(rule)
		if err != nil {
			return err
		}
	case "not_empty":
		err := p.expectNoArgRule(rule)
		if err != nil {
			return err
		}
	case "gt", "lt", "gte", "lte", "len_gt", "len_gte", "len_lt", "len_lte":
		err := p.expectDecimalRule(rule)
		if err != nil {
			return err
		}
	case "custom":
		err := p.expectIdentRule(rule)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected rule %s", rule)
	}

	token = p.Peek()
	if token.typ == COMMA {
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
	default:
		p.AddRule(Rule{
			name:  p.name,
			rule:  rule,
			depth: p.depth,
		})
	}
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
		rule:  rule,
		param: token.value,
		depth: p.depth,
	})
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
		rule:  rule,
		param: token.value,
		depth: p.depth,
	})
	return nil
}

func (p *Parser) expectToken(expected TokenType) (Token, error) {
	token := p.Consume()
	if token.typ != expected {
		return Token{}, fmt.Errorf("unexpected token %s expected %s", token.typ.String(), expected.String())
	}
	return token, nil
}
