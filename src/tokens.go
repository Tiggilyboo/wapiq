package wapiq

type Token int

const (
	// Special
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT        // x
	IDENT_ACTION // x FOR ...
	IDENT_MAP    // ... FOR x
	IDENT_WHERE  // x `x`
	IDENT_NAME   // "x"
	IDENT_VALUE  // `x`

	// Scope
	OBJECT_OPEN  // {
	OBJECT_CLOSE // }
	ARRAY_OPEN   // [
	ARRAY_CLOSE  // ]

	QUOTE_NAME  // "
	QUOTE_VALUE // `

	// State
	DONE // ;

	// Keywords
	API      // API
	API_PATH // path
	API_ARGS // args

	GET       // GET
	POST      // POST
	ACT_PATH  // path
	ACT_ARGS  // args
	ACT_TYPE  // type
	ACT_HEAD  // head
	ACT_QUERY // query
	ACT_BODY  // body

	MAP   // MAP
	QUERY // QUERY
	FOR   // FOR
	WHERE // WHERE
)

func isWS(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNum(r rune) bool {
	return (r >= '0' && r <= '9')
}

func isSym(r rune) bool {
	return r == '!' || (r >= '#' && r <= '&') || (r >= '*' && r <= '/') || (r >= ':' && r <= '@') || (r >= '[' && r <= '\'') || (r >= '{' && r <= '~')
}

const eof rune = rune(0)
