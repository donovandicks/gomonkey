package ast

import (
	"fmt"
	"strings"

	"github.com/donovandicks/gomonkey/internal/token"
)

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

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func NewIntegerLiteral(value int64) *IntegerLiteral {
	return &IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: fmt.Sprint(value)},
		Value: value,
	}
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: value},
		Value: value,
	}
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
func NewBoolean(value bool) *Boolean {
	var tt token.TokenType
	if value {
		tt = token.TRUE
	} else {
		tt = token.FALSE
	}

	return &Boolean{
		Token: token.Token{Type: tt, Literal: fmt.Sprint(value)},
		Value: value,
	}
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out strings.Builder

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out strings.Builder

	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out strings.Builder

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type AssignmentExpression struct {
	Token token.Token // identifier token
	Left  *Identifier
	Right Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out strings.Builder

	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Right.String())

	return out.String()
}

type ListLiteral struct {
	Token token.Token // the '[' token
	Elems []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }
func (ll *ListLiteral) String() string {
	var out strings.Builder

	es := []string{}
	for _, expr := range ll.Elems {
		es = append(es, expr.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")

	return out.String()
}

type MapLiteral struct {
	Token   token.Token // the '{' token
	Entries map[Expression]Expression
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MapLiteral) String() string {
	var out strings.Builder

	kvs := []string{}
	for k, v := range ml.Entries {
		kvs = append(kvs, fmt.Sprintf("%s:%s", k.String(), v.String()))
	}

	out.WriteString(fmt.Sprintf("{%s}", strings.Join(kvs, ", ")))
	return out.String()
}

type IndexExpression struct {
	Token token.Token // the '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
