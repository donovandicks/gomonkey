package parser

import (
	"github.com/donovandicks/gomonkey/internal/token"
)

type (
	PrecedenceTable    map[token.TokenType]OperatorPrecedence
	OperatorPrecedence int
)

const (
	_ OperatorPrecedence = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var Precedence PrecedenceTable = PrecedenceTable{
	token.EQ:     EQUALS,
	token.NE:     EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.FSLASH: PRODUCT,
	token.STAR:   PRODUCT,
	token.LPAREN: CALL,
}
