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
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.FSLASH, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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

func (p *Parser) peekPrecedence() OperatorPrecedence {
	if p, ok := Precedence[p.nextToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() OperatorPrecedence {
	if p, ok := Precedence[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currToken.Type == token.TRUE}
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

	p.readToken() // Read assignment and move on
	p.readToken() // Read the beginning of the expression

	stmt.Value = p.parseExpression(LOWEST)

	for p.currToken.Type != token.SEMICOLON {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.readToken()

	stmt.Value = p.parseExpression(LOWEST)

	for p.currToken.Type != token.SEMICOLON {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	p.readToken()

	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}

	prec := p.currPrecedence()
	p.readToken()
	expr.Right = p.parseExpression(prec)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.readToken() // advance past the '('

	expr := p.parseExpression(LOWEST)

	p.readToken() // advance to the next token after the expression

	if p.currToken.Type != token.RPAREN {
		// expression was parsed but group did not close
		p.addError(ErrMissingCloser{expected: ")"})
		return nil
	}

	return expr
}

func (p *Parser) parseExpression(precedence OperatorPrecedence) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.addError(ErrNoPrefixParser{operator: p.currToken.Literal})
		return nil
	}

	leftExp := prefix()
	for p.nextToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExp
		}

		p.readToken()

		leftExp = infix(leftExp)
	}

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
