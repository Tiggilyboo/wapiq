package parser

import (
	"fmt"
)

type Action struct {
	Token      Token    `json:"t"`
	Identifier string   `json:"i"`
	Value      string   `json:"v"`
	Actions    []Action `json:"a"`
}

func (a *Action) Has(t []Token, v string) bool {
	if a.Token.Exists(t) {
		switch a.Token {
		case T_IDENT_VALUE:
			if a.Value == v {
				return true
			}
		case T_IDENT_NAME:
			if a.Identifier == v {
				return true
			}
		case T_IDENT_PAIR:
			if a.Identifier == v || a.Value == v {
				return true
			}
		}
	}
	for _, aa := range a.Actions {
		return aa.Has(t, v)
	}
	return false
}

func (p *Parser) parseMap(a Action) ([]Action, error) {
	var t Token
	var l string
	var MapTokens = [...]Token{T_IDENT_NAME, T_IDENT_VALUE, T_IDENT_PAIR}

	t, l = p.scanIgnoreWS()
	if t != T_OBJECT_OPEN {
		return nil, fmt.Errorf("Found %q, expected { after MAP action identifier.", l)
	}
	p.unscan()

	acts, err := p.parseObjectBody()
	if err != nil {
		return nil, fmt.Errorf("Error while parsing MAP action %q: %s", a.Identifier, err.Error())
	}
	for _, as := range acts {
		if !as.Has(MapTokens[:], as.Identifier) {
			return nil, fmt.Errorf("Found %q, MAP action %q does not have this query, body or head parameter.", as.Identifier, a.Identifier)
		}
	}

	return acts, nil
}

func (p *Parser) parseAction(maps ParseMap) (*Action, error) {
	var a Action
	var t Token
	var l string
	var ApiActions = [...]Token{T_ACT_ARGS, T_ACT_PATH}
	var HttpActions = [...]Token{T_ACT_PATH, T_ACT_TYPE, T_ACT_BODY, T_ACT_HEAD, T_ACT_QUERY}

	t, l = p.scanIgnoreWS()
	if t != T_IDENT_NAME {
		return nil, fmt.Errorf("Found %q, expected name identifier.", l)
	}
	a.Identifier = l

	cmdToken, cmds := p.scanIgnoreWS()
	if cmdToken != T_API && cmdToken != T_GET && cmdToken != T_POST && cmdToken != T_MAP {
		return nil, fmt.Errorf("Found %q, expected API, GET, POST, MAP command type.", cmds)
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
		availActions = []Token{}
	}

loop:
	for {
		t, l = p.scanIgnoreWS()
		if t == T_OBJECT_CLOSE {
			break loop
		}
		p.unscan()

		actionToken, actionString := p.scanIgnoreWS()
		if len(availActions) > 0 && Token.Exists(actionToken, availActions) == false {
			return nil, fmt.Errorf("Found %q, not an available action for %s command", actionString, cmds)
		} else if len(availActions) == 0 {
			mapAction, mapExists := maps[actionString]
			if !mapExists {
				return nil, fmt.Errorf("Found %q, not an available %s action", actionString, cmds)
			}
			mapBody, err := p.parseMap(mapAction)
			if err != nil {
				return nil, fmt.Errorf("Error parsing %s command's value, expected quoted value: %s", actionString, err.Error())
			}
			a.Actions = append(a.Actions, Action{
				Token:      actionToken,
				Identifier: actionString,
				Actions:    mapBody,
			})
		} else {

			switch actionToken {
			// IDENT_OBJECT
			case T_ACT_ARGS:
				objectBody, err := p.parseObjectBody()
				if err != nil {
					return nil, fmt.Errorf("Error parsing %s command's %s value, expected object: %s", cmds, actionString, err.Error())
				}
				a.Actions = append(a.Actions, Action{
					Token:      actionToken,
					Identifier: actionString,
					Actions:    objectBody,
				})
				// IDENT_ARRAY
			case T_ACT_BODY, T_ACT_HEAD, T_ACT_QUERY:
				arrayBody, err := p.parseArrayBody()
				if err != nil {
					return nil, fmt.Errorf("Error parsing %s command's %s value, expected array: %s", cmds, actionString, err.Error())
				}
				a.Actions = append(a.Actions, Action{
					Token:      actionToken,
					Identifier: actionString,
					Actions:    arrayBody,
				})
				// IDENT_VALUE
			case T_ACT_PATH, T_ACT_TYPE:
				t, l = p.scanIgnoreWS()
				if t != T_IDENT_VALUE {
					return nil, fmt.Errorf("Found %q, Error parsing %s command's %s value, expected quoted value.", l, cmds, actionString)
				}
				a.Actions = append(a.Actions, Action{
					Token:      actionToken,
					Identifier: actionString,
					Value:      l,
				})
			}
		}
	}

	return &a, nil
}
