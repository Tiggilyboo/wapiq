package parser

import "strings"

type Token int

const (
	// Special
	T_ILLEGAL Token = iota
	T_EOF
	T_WS
	T_COMMENT

	// Literals
	T_IDENT        // x
	T_IDENT_WHERE  // x `x`
	T_IDENT_NAME   // "x"
	T_IDENT_VALUE  // `x`
	T_IDENT_PAIR   // "x" `y`
	T_IDENT_ARRAY  // x [ ... ]
	T_IDENT_OBJECT // x { ... }

	// Scope
	T_OBJECT_OPEN  // {
	T_OBJECT_CLOSE // }
	T_ARRAY_OPEN   // [
	T_ARRAY_CLOSE  // ]

	T_QUOTE_NAME  // "
	T_QUOTE_VALUE // `

	// State
	T_DONE // ;

	// Keywords
	T_API // API

	T_GET       // GET
	T_POST      // POST
	T_ACT_PATH  // path
	T_ACT_ARGS  // args
	T_ACT_TYPE  // type
	T_ACT_HEAD  // head
	T_ACT_QUERY // query
	T_ACT_BODY  // body

	T_MAP   // MAP
	T_QUERY // QUERY
	T_FOR   // FOR
	T_WHERE // WHERE
)

func (t Token) Exists(set []Token) bool {
	for _, v := range set {
		if t == v {
			return true
		}
	}
	return false
}

func isWS(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '	'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNum(r rune) bool {
	return (r >= '0' && r <= '9')
}

func isSym(r rune) bool {
	return strings.ContainsRune("~!@#$%^&*()-_=+{}[];:'\",.<>/?\\|", r)
}

const eof rune = rune(0)
