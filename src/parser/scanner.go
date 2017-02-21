package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	r, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) scanComment() (t Token, l string) {
	var b bytes.Buffer
	b.WriteRune(s.read())
	s.unread()

	if r := s.read(); r != '#' {
		s.unread()
		return T_COMMENT, b.String()
	}

	for {
		if r := s.read(); r == eof || r == '\n' {
			break
		} else {
			b.WriteRune(r)
		}
	}

	return T_COMMENT, b.String()
}

func (s *Scanner) scanWS() (t Token, l string) {
	var b bytes.Buffer
	b.WriteRune(s.read())

	for {
		if r := s.read(); r == eof {
			break
		} else if !isWS(r) {
			s.unread()
			break
		} else {
			b.WriteRune(r)
		}
	}

	return T_WS, b.String()
}

func (s *Scanner) scanIdent() (t Token, l string) {
	var b bytes.Buffer
	b.WriteRune(s.read())

	for {
		if r := s.read(); r == eof {
			break
		} else if !isAlpha(r) && !isNum(r) && !isSym(r) {
			s.unread()
			break
		} else {
			b.WriteRune(r)
		}
	}

	// Keywords
	bs := b.String()
	switch strings.ToUpper(bs) {
	case "API":
		return T_API, bs
	case "GET":
		return T_GET, bs
	case "POST":
		return T_POST, bs
	case "MAP":
		return T_MAP, bs
	case "QUERY":
		return T_QUERY, bs
	case "FOR":
		return T_FOR, bs
	case "WHERE":
		return T_WHERE, bs
	}
	return T_IDENT, bs
}

func (s *Scanner) Scan() (t Token, l string) {
	r := s.read()
	if isWS(r) {
		s.unread()
		return s.scanWS()
	} else if r == '#' {
		s.unread()
		return s.scanComment()
	} else if isAlpha(r) || isNum(r) {
		s.unread()
		return s.scanIdent()
	}

	switch r {
	case eof:
		return T_EOF, ""
	case '"':
		return T_QUOTE_NAME, string(r)
	case '`':
		return T_QUOTE_VALUE, string(r)
	case '{':
		return T_OBJECT_OPEN, string(r)
	case '}':
		return T_OBJECT_CLOSE, string(r)
	case '[':
		return T_ARRAY_OPEN, string(r)
	case ']':
		return T_ARRAY_CLOSE, string(r)
	case ';':
		return T_DONE, string(r)
	default:
		return T_ILLEGAL, string(r)
	}
}
