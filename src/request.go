package wapiq

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type RequestType int

const (
	RequestTypeGET RequestType = iota
	RequestTypePOST
)

type RequestFormat int

const (
	RequestFormatJSON RequestFormat = iota
	RequestFormatXML
)

type Request struct {
	Name   string
	Path   string
	Type   RequestType
	Format RequestFormat
	Head   http.Header
	Query  url.Values
	Body   url.Values
}

type Response struct {
	StatusCode int
	Status     string
	Response   bytes.Buffer
}

func (r *Request) URL(a *API) (*url.URL, error) {
	u, err := url.Parse(a.Path + r.Path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = r.Query.Encode()
	return u, nil
}

func (r *Request) fillEmptyArgs(a *API) {
	var vs string
	var e bool
	for k, v := range *a.Args {
		vs = v.(string)
		_, e = r.Query[k]
		if e {
			r.Query.Set(k, vs)
			continue
		}
		_, e = r.Head[k]
		if e {
			r.Head.Set(k, vs)
		}
		_, e = r.Body[k]
		if e {
			r.Body.Set(k, vs)
		}
		p := "{" + k + "}"
		if strings.Contains(r.Path, p) {
			r.Path = strings.Replace(r.Path, p, vs, 1)
		}
	}
	for k, v := range r.Query {
		if len(v) == 0 || v[0] == "" {
			r.Query.Del(k)
		}
	}
}

func (r *Request) Invoke(a *API) (*Response, error) {
	var req *http.Request
	var u *url.URL
	var body io.Reader
	var rb bytes.Buffer
	var err error

	r.fillEmptyArgs(a)

	if len(r.Body) > 0 {
		body = bytes.NewReader([]byte(r.Body.Encode()))
	}
	u, err = r.URL(a)
	if err != nil {
		return nil, err
	}
	switch r.Type {
	case RequestTypeGET:
		req, err = http.NewRequest("GET", u.String(), body)
	case RequestTypePOST:
		req, err = http.NewRequest("POST", u.String(), body)
	}
	if err != nil {
		return nil, err
	}
	if r.Head != nil {
		req.Header = r.Head
	}

	var resp *http.Response
	c := http.Client{}
	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	rw := bufio.NewWriter(&rb)
	_, err = io.Copy(rw, resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Response:   rb,
	}, nil
}
