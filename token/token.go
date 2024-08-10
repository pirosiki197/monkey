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

	IDENT // IDENT
	INT   // INT

	ASSIGN // =
	PLUS   // +

	COMMA     // ,
	SEMICOLON // ;

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	FUNCTION // FUNCTION
	LET      // LET
)
