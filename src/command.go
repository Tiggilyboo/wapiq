package wapiq

import "fmt"

type Command struct {
	Token   Token
	Actions []Action
	State   State
}

type Action struct {
	Token      Token
	Identifier string
	Value      string
}

type State int

const (
	STATE_NONE State = iota
	STATE_ARRAY_PARENT
	STATE_ARRAY_VALUE
	STATE_OBJECT_PARENT
	STATE_OBJECT_VALUE
	STATE_ILLEGAL
)

func (p *Parser) parseCommand() (*Command, error) {
	var t Token
	var l string

	t, l = p.scanQuoted(QUOTE_NAME)
	if t != IDENT_NAME {
		return nil, fmt.Errorf("Found %q, expected name identifier.", l)
	}

	cmd := &Command{
		Token: IDENT_NAME,
		State: STATE_OBJECT_PARENT,
	}

	t, l = p.scanIgnoreWS()
	if t != API && t != GET && t != POST && t != MAP && t != OBJECT_OPEN {
		t, l = p.scanIgnoreWS()
		if t != OBJECT_OPEN {
			return nil, fmt.Errorf("Found %q, expected { after API, GET, POST, or MAP.")
		}
	}
	cmd.Token = t

	for {
		switch t {
		case API:
			t, l = p.scanActionKeyword()
			if t != ACT_PATH && t != ACT_ARGS {
				return nil, fmt.Errorf("Found %q, expected path or args.")
			}
			a := Action{
				Token:      t,
				Identifier: l,
			}
			if t == ACT_PATH {
				t, l = p.scanQuoted(QUOTE_VALUE)
				if t != IDENT_VALUE {
					return nil, fmt.Errorf("Found %q, expected path value.")
				}
				a.Value = l
			} else if t == ACT_ARGS {
				o, oerr := p.parseObject(QUOTE_NAME)
				if oerr != nil {
					return nil, fmt.Errorf("Found %q, expected args object declaration.")
				}
				cmd.Actions = append(cmd.Actions, o.Actions)
			}

		case GET:

		case POST:
		case MAP:
		case OBJECT_OPEN:

		default:
			return nil, fmt.Errorf("Found %q, expected one of API, GET, POST, MAP.")
		}
	}

	return cmd, nil
}
