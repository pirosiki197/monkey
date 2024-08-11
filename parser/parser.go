package parser

import (
	"fmt"

	"github.com/pirosiki197/monkey/ast"
	"github.com/pirosiki197/monkey/lexer"
	"github.com/pirosiki197/monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errs []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: parse LetStatement.Value
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(tok token.TokenType) bool {
	if p.peekToken.Type == tok {
		p.nextToken()
		return true
	} else {
		p.peekError(tok)
		return false
	}
}

func (p *Parser) peekError(tok token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tok, p.peekToken.Type)
	p.errs = append(p.errs, msg)
}

func (p *Parser) Errors() []string {
	return p.errs
}
