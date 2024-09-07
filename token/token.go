package token

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
}

//go:generate stringer -type TokenType -linecomment token.go

const (
	_       TokenType = iota
	ILLEGAL           // ILLEGAL
	EOF               // EOF

	IDENT  // IDENT
	INT    // INT
	STRING // STRING

	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /

	EQ     // ==
	NOT_EQ // !=
	LT     // <
	GT     // >
	LT_EQ  // <=
	GT_EQ  // >=

	COMMA     // ,
	SEMICOLON // ;

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	FUNCTION // FUNCTION
	LET      // LET
	TRUE     // TRUE
	FALSE    // FALSE
	IF       // IF
	ELSE     // ELSE
	RETURN   // RETURN
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
