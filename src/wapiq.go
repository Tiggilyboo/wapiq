package wapiq

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strconv"

	"./parser"
)

type WAPIQ struct {
	APIs        map[string]API
	Maps        map[string]Map
	Requests    map[string]Request
	Queries     []Query
	initialized bool
}

func (*WAPIQ) New() *WAPIQ {
	w := &WAPIQ{}
	w.APIs = map[string]API{}
	w.Maps = map[string]Map{}
	w.Requests = map[string]Request{}
	w.initialized = true
	return w
}

func (w *WAPIQ) serializeActions(actions []parser.Action) error {
	var err error

	for _, act := range actions {
		switch act.Token {
		case parser.T_API:
			var api *API
			api, err = api.Serialize(&act)
			if err != nil {
				return err
			}
			w.APIs[api.Name] = *api
		case parser.T_MAP:
			var m *Map
			m, err = m.Serialize(w.APIs, &act)
			if err != nil {
				return err
			}
			w.Maps[m.Name] = *m
		case parser.T_GET, parser.T_POST:
			var req *Request
			req, err = req.Serialize(&act)
			if err != nil {
				return err
			}
			w.Requests[req.Name] = *req
		case parser.T_QUERY:
			var qy *Query
			qr, err := qy.Serialize(&act)
			if err != nil {
				return err
			}
			w.fillQuery(qr, &act)
			w.Queries = append(w.Queries, *qr)
		}
	}

	return nil
}

func (w *WAPIQ) fillQuery(q *Query, a *parser.Action) error {
	if a.Token != parser.T_QUERY {
		return errors.New("Invalid action to fill query, must be QUERY action.")
	}
	r, re := w.Requests[a.Identifier]
	if !re {
		return errors.New("Found %q, query action does not exist.")
	}
	m, me := w.Maps[a.Value]
	if !me {
		return errors.New("Found %q, query map does not exist.")
	}
	q.Map = m
	q.Request = r

	return nil
}

func (w *WAPIQ) Execute(script string) error {
	if !w.initialized {
		return errors.New("WAPIQ has not been initialized correctly, example use: \n\tw := WAPIQ.New()")
	}
	b := []byte(script)
	r := bytes.NewReader(b)
	p := parser.NewParser(r)
	a, err := p.Parse()
	if err != nil {
		return err
	}
	err = w.serializeActions(a)
	if err != nil {
		return err
	}
	return nil
}

func (w *WAPIQ) Query(script string) (*MapResult, error) {
	err := w.Execute(script)
	if err != nil {
		return nil, err
	}
	mr := MapResult{}
	for i, q := range w.Queries {
		mr[strconv.FormatInt(int64(i), 10)], err = q.Invoke()
		if err != nil {
			return nil, err
		}
	}

	return &mr, nil
}

func (w *WAPIQ) Load(file string, query bool) (*MapResult, error) {
	var b []byte
	var err error
	b, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if query {
		return w.Query(string(b))
	} else {
		return nil, w.Execute(string(b))
	}
}
