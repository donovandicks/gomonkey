package parser

import (
	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	nextToken token.Token
	errors    []string
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(e error) {
	p.errors = append(p.errors, e.Error())
}

func (p *Parser) readToken() {
	p.currToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) expectNext(t token.TokenType) bool {
	return p.nextToken.Type == t
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectNext(token.IDENT) {
		p.addError(ErrNextTokenInvalid{expected: token.IDENT, actual: p.nextToken.Type})
		return nil
	}

	p.readToken()

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectNext(token.ASSIGN) {
		p.addError(ErrNextTokenInvalid{expected: token.ASSIGN, actual: p.nextToken.Type})
		return nil
	}

	p.readToken()

	for p.currToken.Type != token.SEMICOLON {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.readToken()

	for p.currToken.Type != token.SEMICOLON {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.NewProgram()

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.readToken()
	}

	return program
}
