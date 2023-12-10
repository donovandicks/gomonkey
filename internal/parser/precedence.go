package parser

type OperatorPrecedence int

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
