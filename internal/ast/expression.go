package ast

import "github.com/donovandicks/gomonkey/internal/token"

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token // IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
