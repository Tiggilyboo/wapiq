package parser

import "fmt"

func (p *Parser) parseInclude(includes ParseMap) (*Action, error) {
	var t Token
	var l string

	t, l = p.scanIgnoreWS()
	if t != T_INCLUDE {
		return nil, fmt.Errorf("Found %q, expected ^ file inclusion token.", l, t)
	}
	_, mapExists := includes[l]
	if mapExists {
		return nil, fmt.Errorf("Found %q, include has already been loaded. Please check for the duplicate!", l)
	}

	a := &Action{
		Token:      T_INCLUDE,
		Identifier: l,
	}

	return a, nil
}
