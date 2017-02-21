package parser

import (
	"fmt"
	"io"
	"strings"
)

type Parser struct {
	s *Scanner
	b struct {
		t Token
		l string
		n int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) scan() (t Token, l string) {
	if p.b.n != 0 {
		p.b.n = 0
		return p.b.t, p.b.l
	}

	t, l = p.s.Scan()
	p.b.t, p.b.l = t, l

	fmt.Printf("parser.scan: %s\n", l)

	return
}

func (p *Parser) unscan() {
	p.b.n = 1
}

func (p *Parser) scanIgnoreWS() (t Token, l string) {
	t, l = p.scan()
	if t == T_WS || t == T_COMMENT {
		t, l = p.scan()
	}
	return
}

func (p *Parser) scanActionKeyword() (t Token, l string) {
	t, l = p.scanIgnoreWS()
	if t != T_IDENT {
		p.unscan()
		return T_IDENT, l
	}

	switch strings.ToLower(l) {
	case "path":
		return T_ACT_PATH, l
	case "args":
		return T_ACT_ARGS, l
	case "type":
		return T_ACT_TYPE, l
	case "head":
		return T_ACT_HEAD, l
	case "query":
		return T_ACT_QUERY, l
	case "body":
		return T_ACT_BODY, l
	default:
		p.unscan()
		return T_IDENT, l
	}
}

func (p *Parser) scanQuoted(quote Token) (t Token, l string) {
	t, l = p.scanIgnoreWS()
	if t != quote {
		return T_ILLEGAL, l
	}

	t, l = p.scanIgnoreWS()
	fmt.Println("scanQuoted: ", t, T_IDENT, l)
	if t != T_IDENT {
		p.unscan()
		return T_ILLEGAL, l
	}
	r := l

	t, l = p.scanIgnoreWS()
	if t != quote {
		p.unscan()
		return T_ILLEGAL, l
	}

	switch quote {
	case T_QUOTE_NAME:
		return T_IDENT_NAME, r
	case T_QUOTE_VALUE:
		return T_IDENT_VALUE, r
	default:
		return T_IDENT, r
	}
}

func (p *Parser) parseArrayBody() ([]Action, error) {
	var t Token
	var l string
	a := []Action{}

	t, l = p.scanIgnoreWS()
	if t != T_ARRAY_OPEN {
		return a, fmt.Errorf("Found %q, expected array [")
	}

	for {
		t, _ = p.scanIgnoreWS()
		if t == T_ARRAY_CLOSE {
			break
		}
		p.unscan()

		t, l = p.scanQuoted(T_QUOTE_VALUE)
		if t != T_IDENT_VALUE {
			return a, fmt.Errorf("Found %q, expected array value")
		}

		a = append(a, Action{
			Token: t,
			Value: l,
		})
	}

	return a, nil
}

func (p *Parser) parseArray(quote Token) (*Action, error) {
	var t Token
	var l string

	if quote == T_QUOTE_NAME {
		t, l = p.scanQuoted(quote)
		if quote == T_QUOTE_NAME && t != T_IDENT_NAME {
			return nil, fmt.Errorf("Found %q, expected quoted name array identifier")
		}
	} else {
		t, l = p.scanIgnoreWS()
		if t != T_IDENT {
			return nil, fmt.Errorf("Found %q, expected array name identifier")
		}
	}
	a := &Action{
		Token:      t,
		Identifier: l,
	}
	arrayBody, err := p.parseArrayBody()
	if err != nil {
		return nil, fmt.Errorf("Eror parsing array body: %s", err.Error())
	}
	a.Actions = arrayBody

	return a, nil
}

func (p *Parser) parseObjectBody() ([]Action, error) {
	var t Token
	var l string

	t, l = p.scanIgnoreWS()
	if t != T_OBJECT_OPEN {
		return nil, fmt.Errorf("Found %q, expected { after object identifier.")
	}

	a := []Action{}
	for {
		t, l = p.scanIgnoreWS()
		if t == T_OBJECT_CLOSE {
			break
		}
		p.unscan()

		t, l = p.scanQuoted(T_QUOTE_NAME)
		if t != T_IDENT_NAME {
			return nil, fmt.Errorf("Found %q, expected name identifier in object.")
		}
		n := l

		t, l = p.scanQuoted(T_QUOTE_VALUE)
		if t != T_IDENT_VALUE {
			return nil, fmt.Errorf("Found %q, expected value identifier after name in object.")
		}

		a = append(a, Action{
			Token:      T_IDENT_PAIR,
			Identifier: n,
			Value:      l,
		})
	}

	return a, nil
}

func (p *Parser) parseObject(quote Token) (*Action, error) {
	var t Token
	var l string

	if quote == T_QUOTE_NAME {
		t, l = p.scanQuoted(quote)
		if quote == T_QUOTE_NAME && t != T_IDENT_NAME {
			return nil, fmt.Errorf("Found %q, expected quoted object name identifier")
		}
	} else {
		t, l = p.scanIgnoreWS()
		if t != T_IDENT {
			return nil, fmt.Errorf("Found %q, expected object name identifier")
		}
	}
	a := &Action{
		Token:      t,
		Identifier: l,
	}
	objectBody, err := p.parseObjectBody()
	if err != nil {
		return nil, fmt.Errorf("Error parsing object body: %s", err.Error())
	}
	a.Actions = objectBody

	return a, nil
}

func (p *Parser) Parse() (*Action, error) {
	t, l := p.scanIgnoreWS()
	switch t {
	case T_DONE:
		return p.Parse()
	case T_QUERY:
		p.unscan()
		return p.parseQuery()
	case T_QUOTE_NAME:
		p.unscan()
		return p.parseAction()
	case T_EOF:
		return nil, nil
	default:
		return nil, fmt.Errorf("Found %q, expected QUERY or quoted identifier.", l)
	}
}
