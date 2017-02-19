package wapiq

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

	return WS, b.String()
}

func (s *Scanner) scanIdent() (t Token, l string) {
	var b bytes.Buffer
	b.WriteRune(s.read())

	for {
		if r := s.read(); r == eof {
			break
		} else if !isLetter(r) && !isNum(r) && !isSym(r) {
			s.unread()
			break
		} else {
			_, _ := b.WriteRune(r)
		}
	}

	// Keywords
	bs := b.String()
	switch strings.ToUpper(bs) {
	case "API":
		return API, bs
	case "GET":
		return GET, bs
	case "POST":
		return POST, bs
	case "MAP":
		return MAP, bs
	case "QUERY":
		return QUERY, bs
	case "FOR":
		return FOR, bs
	case "WHERE":
		return WHERE, bs
	}

	return IDENT, bs
}

func (s *Scanner) Scan() (t Token, l string) {
	r := s.read()
	if isWS(r) {
		s.unread()
		return s.scanWS()
	} else if isAlpha(r) || isNum(r) {
		s.unread()
		return s.scanIdent()
	}

	switch r {
	case eof:
		return EOF, ""
	case '"':
		return QUOTE_NAME, string(r)
	case '`':
		return QUOTE_VALUE, string(r)
	case '{':
		return OBJECT_OPEN, string(r)
	case '}':
		return OBJECT_CLOSE, string(r)
	case '[':
		return ARRAY_OPEN, string(r)
	case ']':
		return ARRAY_CLOSE, string(r)
	case ';':
		return DONE, string(r)
	default:
		return ILLEGAL, string(r)
	}
}
