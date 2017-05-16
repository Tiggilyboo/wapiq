package parser

import (
	"fmt"
	"io"
	"strings"
)

type ParseMap map[string]Action

type Parser struct {
	s *Scanner
	b struct {
		t Token
		l string
		n int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: NewScanner(r),
	}
}

func (p *Parser) scan() (t Token, l string) {
	if p.b.n != 0 {
		p.b.n = 0
		return p.b.t, p.b.l
	}

	t, l = p.s.Scan()
	p.b.t, p.b.l = t, l
	return
}

func (p *Parser) unscan() {
	p.b.n = 1
}

func (p *Parser) scanIgnoreWS() (t Token, l string) {
	t, l = p.scan()
	for {
		if t == T_WS || t == T_COMMENT {
			t, l = p.scan()
		} else {
			break
		}
	}
	return t, l
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
		t, l = p.scanIgnoreWS()
		if t == T_ARRAY_CLOSE {
			break
		} else if t != T_IDENT_VALUE {
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
		t, l = p.scanIgnoreWS()
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
		return nil, fmt.Errorf("Found %q, expected { after object identifier.", l)
	}

	a := []Action{}
	for {
		t, l = p.scanIgnoreWS()
		if t == T_OBJECT_CLOSE {
			break
		} else if t != T_IDENT_NAME {
			return nil, fmt.Errorf("Found %q, expected name identifier in object.", l)
		}
		n := l

		t, l = p.scanIgnoreWS()
		if t != T_IDENT_VALUE && t != T_AT {
			return nil, fmt.Errorf("Found %q, expected value identifier (surrounded by `) or an @ expression after name identifier.")
		}
		if t == T_IDENT_VALUE {
			a = append(a, Action{
				Token:      T_IDENT_PAIR,
				Identifier: n,
				Value:      l,
			})
		} else if t == T_AT {
			indices := []string{}
			for {
				t, l = p.scanIgnoreWS()
				if t == T_IDENT_NAME || t == T_OBJECT_CLOSE {
					p.unscan()
					break
				} else if t != T_IDENT && l != "," {
					return nil, fmt.Errorf("Found %q, expected numerical array indices after @ expression.", l)
				}
				indices = append(indices, l)
			}
			a = append(a, Action{
				Token:      T_IDENT_PAIR_AT,
				Identifier: n,
				Value:      strings.Join(indices, ","),
			})
		}
	}

	return a, nil
}

func (p *Parser) parseObject(quote Token) (*Action, error) {
	var t Token
	var l string

	if quote == T_QUOTE_NAME {
		t, l = p.scanIgnoreWS()
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

func (p *Parser) Parse() ([]Action, error) {
	var err error
	var pa *Action
	maps := ParseMap{}
	queries := ParseMap{}
	apis := ParseMap{}
	includes := ParseMap{}
	a := []Action{}

loop:
	for {
		t, l := p.scanIgnoreWS()
		switch t {
		case T_COMMENT:
			break
		case T_INCLUDE:
			p.unscan()
			pa, err = p.parseInclude(includes)
			includes[pa.Identifier] = *pa
			a = append(a, *pa)
		case T_DONE:
			if err == nil && pa != nil {
				switch pa.Token {
				case T_GET, T_POST:
					maps[pa.Identifier] = *pa
				case T_API:
					apis[pa.Identifier] = *pa
				case T_QUERY:
					fmt.Println(pa)
					queries[pa.Identifier] = *pa
				}
				a = append(a, *pa)
			}
		case T_QUERY:
			p.unscan()
			pa, err = p.parseQuery()
		case T_IDENT_NAME, T_IDENT_VALUE:
			p.unscan()
			pa, err = p.parseAction(maps)
		case T_EOF:
			break loop
		default:
			if err == nil {
				err = fmt.Errorf("Found %q, expected query (/), inclusion (^) or quoted identifier followed by GET, POST or MAP.", l)
			}
			break loop
		}
		if err != nil {
			break loop
		}
	}

	return a, err
}
