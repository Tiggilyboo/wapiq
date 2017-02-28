package wapiq

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"./parser"
)

func serializeArray(a *parser.Action) []string {
	r := make([]string, len(a.Actions))
	for i, c := range a.Actions {
		r[i] = c.Value
	}
	return r
}
func serializeObject(a *parser.Action) *MapResult {
	r := MapResult{}
	for _, c := range a.Actions {
		r[c.Identifier] = c.Value
	}
	return &r
}

func (*API) Serialize(a *parser.Action) (*API, error) {
	if a.Token != parser.T_API {
		return nil, errors.New("Invalid parser token, expected API action.")
	}
	api := API{
		Name: a.Identifier,
	}
	for _, act := range a.Actions {
		switch act.Token {
		case parser.T_ACT_PATH:
			api.Path = act.Value
		case parser.T_ACT_ARGS:
			api.Args = serializeObject(&act)
		}
	}

	return &api, nil
}

func (*Map) Serialize(apis map[string]API, a *parser.Action) (*Map, error) {
	if a.Token != parser.T_MAP {
		return nil, errors.New("Invalid parser token, expected MAP action.")
	}
	if len(a.Actions) == 0 {
		return nil, fmt.Errorf("Map %q contains no actions to map against.", a.Identifier)
	}
	api, ae := apis[a.Value]
	if !ae {
		return nil, fmt.Errorf("Map %q references non-existant API %q.", a.Identifier, a.Value)
	}
	m := Map{
		Name:     a.Identifier,
		API:      &api,
		FieldMap: serializeObject(&a.Actions[0]),
		Request:  a.Actions[0].Identifier,
	}

	return &m, nil
}

func (*Request) Serialize(a *parser.Action) (*Request, error) {
	if a.Token != parser.T_GET && a.Token != parser.T_POST {
		return nil, errors.New("Invalid parser token, expected GET or POST action.")
	}
	t := RequestTypeGET
	if a.Token != parser.T_GET {
		t = RequestTypePOST
	}
	r := Request{
		Name:  a.Identifier,
		Type:  t,
		Head:  http.Header{},
		Body:  url.Values{},
		Query: url.Values{},
	}
	for _, act := range a.Actions {
		switch act.Token {
		case parser.T_ACT_PATH:
			r.Path = act.Value
		case parser.T_ACT_TYPE:
			switch strings.ToLower(act.Value) {
			case "json":
				r.Format = RequestFormatJSON
			case "xml":
				r.Format = RequestFormatXML
			}
		case parser.T_ACT_HEAD:
			h := serializeArray(&act)
			for _, hv := range h {
				r.Head.Add(hv, "")
			}
		case parser.T_ACT_BODY:
			b := serializeArray(&act)
			for _, bv := range b {
				r.Body.Add(bv, "")
			}
		case parser.T_ACT_QUERY:
			q := serializeArray(&act)
			for _, qv := range q {
				r.Query.Add(qv, "")
			}
		}
	}
	return &r, nil
}

func (*Query) Serialize(a *parser.Action) (*Query, error) {
	if a.Token != parser.T_QUERY {
		return nil, errors.New("Invalid parser token, expected QUERY")
	}
	args := serializeObject(a)
	q := Query{
		Args: *args,
	}
	return &q, nil
}
