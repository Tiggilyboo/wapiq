package parser

import "fmt"

type Action struct {
	Token      Token
	Identifier string
	Value      string
	Actions    []Action
}

func (p *Parser) parseAction() (*Action, error) {
	var a Action
	var t Token
	var l string
	var ApiActions = [...]Token{T_ACT_ARGS, T_ACT_PATH}
	var HttpActions = [...]Token{T_ACT_PATH, T_ACT_TYPE, T_ACT_BODY, T_ACT_HEAD, T_ACT_QUERY}
	var MapActions = [...]Token{}

	t, l = p.scanQuoted(T_QUOTE_NAME)
	if t != T_IDENT_NAME {
		return nil, fmt.Errorf("Found %q, expected name identifier.", l)
	}

	fmt.Println("Identifier: ", l)
	a.Identifier = l

	cmdToken, cmds := p.scanIgnoreWS()
	if cmdToken != T_API && cmdToken != T_GET && cmdToken != T_POST && cmdToken != T_MAP {
		return nil, fmt.Errorf("Found %q, expected API, GET, POST, MAP command type.", l)
	}
	a.Token = cmdToken

	t, l = p.scanIgnoreWS()
	if t != T_OBJECT_OPEN {
		return nil, fmt.Errorf("Found %q, expected { for %s command.", l, cmds)
	}
	var availActions []Token
	switch cmdToken {
	case T_GET, T_POST:
		availActions = HttpActions[:]
	case T_API:
		availActions = ApiActions[:]
	case T_MAP:
		availActions = MapActions[:]
	}

loop:
	for {
		t, l = p.scanIgnoreWS()
		if t == T_OBJECT_CLOSE {
			break loop
		}
		p.unscan()

		actionToken, actionString := p.scanActionKeyword()
		if Token.Exists(actionToken, availActions) == false {
			return nil, fmt.Errorf("Found %q, not an available action for %s command", actionString, cmds)
		}

		switch actionToken {
		// IDENT_OBJECT
		case T_ACT_ARGS:
			objectBody, err := p.parseObjectBody()
			if err != nil {
				return nil, fmt.Errorf("Error parsing %s command's %s value, expected object: %s", cmds, actionString, err.Error())
			}
			a.Actions = append(a.Actions, objectBody...)

			// IDENT_ARRAY
		case T_ACT_BODY, T_ACT_HEAD, T_ACT_QUERY:
			arrayBody, err := p.parseArrayBody()
			if err != nil {
				return nil, fmt.Errorf("Error parsing %s command's %s value, expected array: %s", cmds, actionString, err.Error())
			}
			a.Actions = append(a.Actions, arrayBody...)

			// IDENT_VALUE
		case T_ACT_PATH, T_ACT_TYPE:
			t, l = p.scanQuoted(T_QUOTE_VALUE)
			if t != T_IDENT_VALUE {
				return nil, fmt.Errorf("Found %q, Error parsing %s command's %s value, expected quoted value.", l, cmds, actionString)
			}
			a.Value = l
		}
	}

	return &a, nil
}
