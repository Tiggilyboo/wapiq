package parser

import (
	"bufio"
	"bytes"
	"fmt"
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
	esc := false
	b.WriteRune(s.read())

	for {
		if r := s.read(); r == eof {
			break
		} else if !esc && r == '\\' {
			esc = true
		} else if !esc && !isAlpha(r) && !isNum(r) {
			s.unread()
			break
		} else {
			if esc && isSym(r) {
				esc = false
			}
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
	case "FOR":
		return T_FOR, bs
	case "WHERE":
		return T_WHERE, bs
	}
	switch strings.ToLower(bs) {
	case "path":
		return T_ACT_PATH, bs
	case "args":
		return T_ACT_ARGS, bs
	case "type":
		return T_ACT_TYPE, bs
	case "head":
		return T_ACT_HEAD, bs
	case "query":
		return T_ACT_QUERY, bs
	case "body":
		return T_ACT_BODY, bs
	}

	return T_IDENT, bs
}

func (s *Scanner) scanQuoted(q rune) (t Token, l string) {
	var b bytes.Buffer
	if q == '"' {
		t = T_IDENT_NAME
	} else if q == '`' {
		t = T_IDENT_VALUE
	} else {
		return T_ILLEGAL, ""
	}

	if r := s.read(); r == eof || r != q {
		fmt.Println("illegal", string(r), string(q))
		return T_ILLEGAL, ""
	}

	esc := false
	for {
		if r := s.read(); !esc && r == q {
			break
		} else if !esc && r == '\\' {
			esc = true
		} else {
			if esc && isSym(r) {
				esc = false
			}
			b.WriteRune(r)
		}
	}

	return t, b.String()
}

func (s *Scanner) Scan() (t Token, l string) {
	r := s.read()
	if isWS(r) {
		s.unread()
		return s.scanWS()
	} else if r == '#' {
		s.unread()
		return s.scanComment()
	} else if r == '"' || r == '`' {
		s.unread()
		return s.scanQuoted(r)
	} else if isAlpha(r) || isNum(r) || r == '\\' {
		s.unread()
		return s.scanIdent()
	}

	switch r {
	case eof:
		return T_EOF, ""
	case '/':
		return T_QUERY, string(r)
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
