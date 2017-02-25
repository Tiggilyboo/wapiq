package main

import (
	"bufio"
	"errors"

	"github.com/Tiggilyboo/gabs"
)

type Map struct {
	Name     string
	Request  string
	FieldMap *MapResult
}

type MapResult map[string]interface{}

func (m *Map) mapJson(r *Response) (*MapResult, error) {
	re := bufio.NewReader(&r.Response)
	gc, err := gabs.ParseJSONBuffer(re)
	if err != nil {
		return nil, err
	}

	mr := MapResult{}
	if m.FieldMap == nil {
		return nil, errors.New("Map contains no fields to map.")
	}
	for k, v := range *m.FieldMap {
		if !gc.Exists(v.(string)) {
			mr[k] = nil
		} else {
			mr[k] = gc.Path(v.(string)).Data()
		}
	}

	return &mr, nil
}

func (m *Map) Invoke(a *API, r *Request) (*MapResult, error) {
	resp, err := r.Invoke(*a)
	if err != nil {
		return nil, err
	}

	switch r.Format {
	case RequestFormatJSON:
		return m.mapJson(resp)
	case RequestFormatXML:
		return nil, errors.New("Unsupported request format type XML, JSON requests only supported at this time! Sorry!")
	}

	return nil, errors.New("Unsupported request format type.")
}
