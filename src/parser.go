package wapiq

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
		return b.b.t, p.b.l
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
	if t == WS {
		t, l = p.scan()
	}
	return
}

func (p *Parser) scanActionKeyword() (t Token, l string) {
	t, l = p.scanIgnoreWS()
	if t != IDENT {
		p.unscan()
		return IDENT, l
	}

	switch strings.ToLower(l) {
	case "path":
		return ACT_PATH, l
	case "args":
		return ACT_ARGS, l
	case "type":
		return ACT_TYPE, l
	case "head":
		return ACT_HEAD, l
	case "query":
		return ACT_QUERY, l
	case "body":
		return ACT_BODY, l
	default:
		p.unscan()
		return IDENT, l
	}
}

func (p *Parser) scanQuoted(quote Token) (t Token, l string) {
	t, l = p.scanIgnoreWS()
	if t != quote {
		p.unscan()
		return ILLEGAL, l
	}

	t, l = p.scanIgnoreWS()
	if t != IDENT {
		p.unscan()
		return ILLEGAL, l
	}
	r := l

	t, l = p.scanIgnoreWS()
	if t != quote {
		p.unscan()
		return ILLEGAL, l
	}

	switch quote {
	case QUOTE_NAME:
		return IDENT_NAME, r
	case QUOTE_VALUE:
		return IDENT_VALUE, r
	default:
		return IDENT, r
	}
}

func (p *Parser) parseObject(quote Token) (*Command, error) {
	var t Token
	var l string

	if quote == QUOTE_NAME || quote == QUOTE_VALUE {
		t, l = p.scanQuoted()
		if quote == QUOTE_NAME && t != IDENT_NAME {
			return nil, fmt.Errorf("Found %q, expected quoted name identifier")
		}
		if quote == QUOTE_VALUE && t != IDENT_VALUE {
			return nil, fmt.Errorf("Found %q, expected quoted value identifier")
		}
	} else {
		t, l = p.scanIgnoreWS()
		if t != IDENT {
			return nil, fmt.Errorf("Found %q, expected name identifier")
		}
	}
	a := Action{
		Token:      t,
		Identifier: l,
	}
	cmd := &Command{
		Token: t,
		State: STATE_OBJECT_PARENT,
	}
	cmd.Actions = append(cmd.Actions, a)

	t, l = p.scanIgnoreWS()
	if t != OBJECT_OPEN {
		return nil, fmt.Errorf("Found %q, expected { after object identifier.")
	}
	p.unscan()

	// Object body
	for {
		t, l = p.scanIgnoreWS()
		if t != OBJECT_CLOSE && t != QUOTE_NAME {
			return nil, fmt.Errorf("Found %q, expected name identifier or } after {.")
		}
		p.unscan()
		if t == OBJECT_CLOSE {
			break
		}

		t, l = p.scanQuoted(QUOTE_NAME)
		if t != IDENT_NAME {
			return nil, fmt.Errorf("Found %q, expected name identifier in object.")
		}
		n := l

		t, l = p.scanQuoted(QUOTE_VALUE)
		if t != IDENT_VALUE {
			return nil, fmt.Errorf("Found %q, expected value identifier after name in object.")
		}

		cmd.Actions = append(cmd.Actions, Action{
			Token:      IDENT,
			Identifier: n,
			Value:      l,
		})
	}

	return cmd, nil
}

func (p *Parser) Parse() (*Command, error) {
	t, l := p.scanIgnoreWS()
	switch t {
	case DONE:
		return p.Parse()
	case QUERY:
		p.unscan()
		return p.parseQuery()
	case QUOTE_NAME:
		p.unscan()
		return p.parseCommand()
	case EOF:
		return nil, nil
	default:
		return nil, fmt.Errorf("Found %q, expected QUERY or quoted identifier.", l)
	}
}
