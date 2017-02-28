package wapiq

import (
	"bufio"
	"errors"
	"reflect"

	"github.com/Tiggilyboo/gabs"
)

type Map struct {
	Name     string
	Request  string
	API      *API
	FieldMap *MapResult
}

type MapResult map[string]interface{}

func (mr *MapResult) JSON() (string, error) {
	gc, err := gabs.Consume(mr)
	if err != nil {
		return "", err
	}
	return gc.String(), nil
}

func (m *Map) mapJson(r *Response) ([]MapResult, error) {
	re := bufio.NewReader(&r.Response)
	gc, err := gabs.ParseJSONBuffer(re)
	if err != nil {
		return nil, err
	}

	mra := []MapResult{}
	if m.FieldMap == nil {
		return nil, errors.New("Map contains no fields to map.")
	}
	for k, v := range *m.FieldMap {
		path := v.(string)
		if gc.ExistsP(path) {
			ctr := gc.Path(path).Data()
			typectr := reflect.TypeOf(ctr)
			switch typectr.Kind() {
			case reflect.Slice:
				valctr := reflect.ValueOf(ctr)
				if len(mra) != valctr.Len() {
					more := make([]MapResult, valctr.Len()-len(mra))
					for mi, _ := range more {
						more[mi] = MapResult{}
					}
					mra = append(mra, more...)
				}
				for i := 0; i < valctr.Len(); i++ {
					mra[i][k] = valctr.Index(i).Interface()
				}
			default:
				if len(mra) == 0 {
					mra = append([]MapResult{}, MapResult{})
				}
				mra[0][k] = ctr
			}
		}
	}

	return mra, nil
}

func (m *Map) Invoke(r *Request) ([]MapResult, error) {
	resp, err := r.Invoke(m.API)
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
