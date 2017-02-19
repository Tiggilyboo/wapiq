package wapiq

import "fmt"

func (p *Parser) parseQuery() (*Command, error) {
	cmd := &Command{
		Token:   QUERY,
		State:   STATE_NONE,
		Actions: make([]Action, 3),
	}
	var t Token
	var l string

	t, l = p.scanIgnoreWS()
	if t != QUERY {
		return nil, fmt.Errorf("Found %q, expected query identifier.", l)
	}

	t, l = p.scanIgnoreWS()
	if t != IDENT {
		return nil, fmt.Errorf("Found %q, expected action identifier.", l)
	}
	cmd.Actions[0] = Action{
		Identifier: l,
		Token:      IDENT_ACTION,
	}

	t, l = p.scanIgnoreWS()
	if t != FOR {
		return nil, fmt.Errorf("Found %q, expected FOR after action identifier.", l)
	}

	t, l = p.scanIgnoreWS()
	if t != IDENT {
		return nil, fmt.Errorf("Found %q, expected map identifier.", l)
	}
	cmd.Actions[1] = Action{
		Identifier: l,
		Token:      IDENT_MAP,
	}

	t, l = p.scanIgnoreWS()
	if t != WHERE && t != DONE {
		return nil, fmt.Errorf("Found %q, expected WHERE or ';'.", l)
	}
	if t == WHERE {
		t, l = p.scanIgnoreWS()
		if t != IDENT {
			return nil, fmt.Errorf("Found %q, expected identifier for query, head, body argument.", l)
		} else {
			p.unscan()
			i := 0

			for {
				t, l = p.scanIgnoreWS()
				if i == 0 && t != IDENT {
					return nil, fmt.Errorf("Found %q, expected name identifier after WHERE or last value in WHERE clause.", l)
				} else if i > 0 {
					break
				}
				a := Action{
					Token:      IDENT_WHERE,
					Identifier: l,
				}

				t, l := p.scanQuoted(QUOTE_VALUE)
				if t != IDENT_VALUE {
					return nil, fmt.Errorf("Found %q, expected associated value after name identifier in WHERE clause.", l)
				}

				cmd.Actions = append(cmd.Actions, a)
				i++
			}
		}
	}

	// DONE
	return cmd, nil
}
