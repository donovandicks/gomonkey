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
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACK, p.parseListLiteral)
	p.registerPrefix(token.LBRACE, p.parseMapLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.FSLASH, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NE, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.DOT, p.parseCallExpression)
	p.registerInfix(token.LBRACK, p.parseIndexExpression)

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

func (p *Parser) parseIdentifier() ast.Expression { return ast.NewIdentifier(p.currToken.Literal) }

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

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseListElements(end token.TokenType) []ast.Expression {
	var exprs []ast.Expression

	if p.nextToken.Type == end {
		p.readToken()
		return exprs
	}

	p.readToken()

	exprs = append(exprs, p.parseExpression(LOWEST))

	for p.nextToken.Type == token.COMMA {
		p.readToken() // advance to ','
		p.readToken() // advance to next expression
		exprs = append(exprs, p.parseExpression(LOWEST))
	}

	if !p.expectNext(end) {
		return nil
	}

	p.readToken()
	return exprs
}

func (p *Parser) parseListLiteral() ast.Expression {
	lit := &ast.ListLiteral{Token: p.currToken}
	lit.Elems = p.parseListElements(token.RBRACK)
	return lit
}

func (p *Parser) parseMapLiteral() ast.Expression {
	m := &ast.MapLiteral{Token: p.currToken}
	m.Entries = make(map[ast.Expression]ast.Expression)

	for !p.expectNext(token.RBRACE) {
		p.readToken() // advance to the key expression

		key := p.parseExpression(LOWEST)
		if !p.expectNext(token.COLON) {
			p.addError(ErrNextTokenInvalid{expected: token.COMMA, actual: p.nextToken.Type})
			return nil
		}

		p.readToken() // advance to the colon
		p.readToken() // advnace to the value expression
		val := p.parseExpression(LOWEST)

		m.Entries[key] = val
		if !p.expectNext(token.RBRACE) && !p.expectNext(token.COMMA) {
			p.addError(ErrNextTokenInvalid{expected: token.RBRACE, actual: p.nextToken.Type})
			return nil
		}

		if p.expectNext(token.COMMA) {
			p.readToken() // advance to the comma
		}
	}

	if !p.expectNext(token.RBRACE) {
		return nil
	}

	p.readToken()
	return m
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currToken.Type == token.TRUE}
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

func (p *Parser) parseAssignmentExpression(operator token.Token, left, right ast.Expression) ast.Expression {
	return &ast.AssignmentExpression{
		Token: operator,
		Left:  left,
		Right: right,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	operator := p.currToken
	prec := p.currPrecedence()

	p.readToken()
	right := p.parseExpression(prec)

	if operator.Type == token.ASSIGN {
		return p.parseAssignmentExpression(operator, left, right)
	}

	return &ast.InfixExpression{
		Token:    operator,
		Operator: operator.Literal,
		Left:     left,
		Right:    right,
	}
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

func (p *Parser) parseWhileStatement() ast.Statement {
	stmt := &ast.WhileStatement{Token: p.currToken}

	if !p.expectNext(token.LPAREN) {
		p.addError(ErrMissingOpener{expected: "("})
		return nil
	}

	p.readToken() // advance to '('
	p.readToken() // advance to condition

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectNext(token.RPAREN) {
		p.addError(ErrMissingCloser{expected: ")"})
		return nil
	}

	p.readToken() // advance to closing paren

	if !p.expectNext(token.LBRACE) {
		p.addError(ErrMissingOpener{expected: "{"})
		return nil
	}

	p.readToken()

	stmt.Block = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currToken}

	if !p.expectNext(token.LPAREN) {
		p.addError(ErrMissingOpener{expected: "("})
		return nil
	}

	p.readToken() // advance to the '('
	p.readToken() // advance to the expression after '('

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectNext(token.RPAREN) {
		p.addError(ErrMissingCloser{expected: ")"})
		return nil
	}

	p.readToken() // advance to the ')'

	if !p.expectNext(token.LBRACE) {
		p.addError(ErrMissingOpener{expected: "{"})
		return nil
	}

	p.readToken() // advance to the '{'

	expr.Consequence = p.parseBlockStatement()

	if p.expectNext(token.ELSE) {
		p.readToken() // advance to the 'else'

		if !p.expectNext(token.LBRACE) {
			p.addError(ErrMissingOpener{expected: "{"})
			return nil
		}

		p.readToken()

		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var idents []*ast.Identifier

	if p.expectNext(token.RPAREN) {
		p.readToken() // advance to the closing ')'
		return idents
	}

	p.readToken()

	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	idents = append(idents, ident)

	for p.expectNext(token.COMMA) {
		p.readToken() // advance to the ','
		p.readToken() // advance to the next ident

		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		idents = append(idents, ident)
	}

	if !p.expectNext(token.RPAREN) {
		p.addError(ErrMissingCloser{expected: ")"})
		return nil
	}

	p.readToken()
	return idents
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectNext(token.LPAREN) {
		p.addError(ErrMissingOpener{expected: "("})
		return nil
	}

	p.readToken() // advance to the '('

	fn.Parameters = p.parseFunctionParameters()

	// Currently on the ')' if one was present
	if !p.expectNext(token.LBRACE) {
		p.addError(ErrMissingOpener{expected: "{"})
		return nil
	}

	p.readToken() // advance to the '{'
	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseCallExpression(callable ast.Expression) ast.Expression {
	curr := p.currToken
	switch curr.Type {
	case token.LPAREN:
		expr := &ast.CallExpression{Token: curr, Function: callable}
		expr.Arguments = p.parseListElements(token.RPAREN)
		return expr
	case token.DOT:
		if !p.expectNext(token.IDENT) {
			p.addError(ErrNextTokenInvalid{expected: token.IDENT})
			return nil
		}

		p.readToken()

		property := p.parseExpression(CALL)

		return &ast.GetExpression{
			Token: curr,
			Left:  callable,
			Right: property,
		}
	default:
		p.addError(ErrParseError{expected: "callable", actual: callable.TokenLiteral()})
		return nil
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.currToken, Left: left}

	p.readToken()

	expr.Index = p.parseExpression(LOWEST)

	if !p.expectNext(token.RBRACK) {
		return nil
	}

	p.readToken()
	return expr
}

func (p *Parser) parseExpression(precedence OperatorPrecedence) ast.Expression {
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.addError(ErrNoPrefixParser{operator: p.currToken.Literal})
		return nil
	}

	leftExp := prefix()
	for !p.expectNext(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExp
		}

		p.readToken()

		leftExp = infix(leftExp)
	}

	return leftExp
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

	if p.expectNext(token.SEMICOLON) {
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

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.expectNext(token.SEMICOLON) {
		p.readToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	p.readToken() // advance past the opening '{'

	block := ast.NewBlock(p.currToken)

	for p.currToken.Type != token.RBRACE && p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.readToken()
	}

	return block
}

func (p *Parser) parseFunctionStatement() ast.Statement {
	fn := &ast.FunctionStatement{Token: token.NewKeyword("fn")}

	// expect the name of the function
	if !p.expectNext(token.IDENT) {
		p.addError(ErrNextTokenInvalid{expected: token.IDENT, actual: p.nextToken.Type})
		return nil
	}

	p.readToken() // advance to the function name

	name, ok := p.parseIdentifier().(*ast.Identifier)
	if !ok {
		p.addError(ErrParseError{expected: token.IDENT, actual: p.currToken.Literal})
		return nil
	}

	fn.Name = name

	p.readToken() // advance to the '(' around the parameters

	fn.Parameters = p.parseFunctionParameters()

	// expect to begin the function body
	if !p.expectNext(token.LBRACE) {
		p.addError(ErrMissingOpener{expected: "{"})
		return nil
	}

	p.readToken() // advance to the opening brace

	fn.Body = p.parseBlockStatement()
	return fn
}

func (p *Parser) parseClassStatement() ast.Statement {
	cs := &ast.ClassStatement{Token: p.currToken}
	if !p.expectNext(token.IDENT) {
		p.addError(ErrNextTokenInvalid{expected: token.IDENT, actual: p.nextToken.Type})
		return nil
	}

	p.readToken() // advance to the class name

	name, ok := p.parseIdentifier().(*ast.Identifier)
	if !ok {
		p.addError(ErrParseError{expected: token.IDENT, actual: string(name.Token.Type)})
		return nil
	}

	cs.Name = name

	if !p.expectNext(token.LBRACE) {
		p.addError(ErrMissingOpener{expected: "{"})
		return nil
	}

	p.readToken() // consume the opening brace

	for p.expectNext(token.IDENT) {
		f := p.parseFunctionStatement()
		fn, ok := f.(*ast.FunctionStatement)
		if !ok {
			return f
		}

		cs.Methods = append(cs.Methods, fn)
	}

	if !p.expectNext(token.RBRACE) {
		p.addError(ErrMissingCloser{expected: "}"})
		return nil
	}

	p.readToken() // consume the closing brace

	return cs
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.CLASS:
		return p.parseClassStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
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
