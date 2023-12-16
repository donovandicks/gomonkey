package ast

import (
	"fmt"
	"strings"

	"github.com/donovandicks/gomonkey/internal/token"
)

type Statement interface {
	Node
	statementNode()
}

type LetStatement struct {
	Token token.Token // LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out strings.Builder

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token token.Token // the return token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of an expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type WhileStatement struct {
	Token     token.Token // the 'while' token
	Condition Expression
	Block     *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out strings.Builder

	stmts := []string{}
	for _, s := range ws.Block.Statements {
		stmts = append(stmts, s.String())
	}

	out.WriteString("while")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Block.String())

	return out.String()
}

type FunctionStatement struct {
	Token      token.Token // the `fn` token
	Name       *Identifier // the function name identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out strings.Builder

	params := []string{}
	for _, param := range fs.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Name.String())
	out.WriteString(fmt.Sprintf("(%s)", strings.Join(params, ", ")))
	out.WriteString(fs.Body.String())

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func NewBlock(currToken token.Token) *BlockStatement {
	return &BlockStatement{Token: currToken}
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ClassStatement struct {
	Token token.Token // the `class` token
}
