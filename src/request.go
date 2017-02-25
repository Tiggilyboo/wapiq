package main

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/url"
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

func (r *Request) URL(a API) (*url.URL, error) {
	u, err := url.Parse(r.Path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = r.Query.Encode()
	return u, nil
}

func (r *Request) Invoke(a API) (*Response, error) {
	var req *http.Request
	var u *url.URL
	var body io.Reader
	var rb bytes.Buffer
	var err error

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

	req.Header = r.Head
	err = req.ParseForm()
	if err != nil {
		return nil, err
	}

	defer req.Response.Body.Close()
	rw := bufio.NewWriter(&rb)
	_, err = io.Copy(rw, req.Response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: req.Response.StatusCode,
		Status:     req.Response.Status,
		Response:   rb,
	}, nil
}
