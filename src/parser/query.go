package parser

import "fmt"

func (p *Parser) parseQuery(maps ParseMap) (*Action, error) {
	var t Token
	var l string
	a := &Action{
		Token:   T_QUERY,
		Actions: []Action{},
	}

	t, l = p.scanIgnoreWS()
	if t != T_QUERY {
		return nil, fmt.Errorf("Found %q, expected query identifier.", l)
	}

	t, l = p.scanIgnoreWS()
	if t != T_IDENT {
		return nil, fmt.Errorf("Found %q, expected action identifier.", l)
	}
	_, e := maps[l]
	if !e {
		return nil, fmt.Errorf("Found %q, unrecognised action idenitifier.", l)
	}
	a.Identifier = l

	t, l = p.scanIgnoreWS()
	if t != T_FOR {
		return nil, fmt.Errorf("Found %q, expected FOR after action identifier.", l)
	}

	t, l = p.scanIgnoreWS()
	if t != T_IDENT {
		return nil, fmt.Errorf("Found %q, expected map identifier.", l)
	}
	a.Value = l

	t, l = p.scanIgnoreWS()
	if t != T_WHERE && t != T_DONE {
		return nil, fmt.Errorf("Found %q, expected WHERE or ;", l)
	}
	if t != T_WHERE {
		return a, nil
	}

	for {
		t, l = p.scanIgnoreWS()
		p.unscan()
		if t == T_DONE {
			break
		}

		t, l = p.scanIgnoreWS()
		if t != T_IDENT {
			return nil, fmt.Errorf("Found %q, expected WHERE clause identifier.", l)
		}
		w := Action{
			Token:      T_IDENT_WHERE,
			Identifier: l,
		}

		t, l = p.scanIgnoreWS()
		if t != T_IDENT_VALUE {
			return nil, fmt.Errorf("Found %q, expected value after '%s' identifier.", l)
		}
		w.Value = l

		a.Actions = append(a.Actions, w)
	}

	return a, nil
}
