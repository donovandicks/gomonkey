package parser

import (
	"strconv"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/token"
)

type (
	prefixParseFn    func() ast.Expression
	infixParseFn     func(ast.Expression) ast.Expression
	PrefixParseFnMap map[token.TokenType]prefixParseFn
	InfixParseFnMap  map[token.TokenType]infixParseFn
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	nextToken token.Token
	errors    []string

	prefixParseFns PrefixParseFnMap
	infixParseFns  InfixParseFnMap
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		prefixParseFns: make(PrefixParseFnMap),
		infixParseFns:  make(InfixParseFnMap),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	p.readToken()
	p.readToken()

	return p
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		p.addError(ErrParseError{actual: p.currToken.Literal, expected: "integer"})
		return nil
	}

	lit.Value = val
	return lit
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

func (p *Parser) parseExpression(precedence OperatorPrecedence) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextToken.Type == token.SEMICOLON {
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
		return p.parseExpressionStatement()
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
