package parser_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/parser"
	"github.com/donovandicks/gomonkey/internal/token"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestParser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		input        string
		expected     []ast.Statement
		expectedErrs []string
	}{
		{
			name:  "valid identifier",
			input: "foo;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Expression: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "foo"},
						Value: "foo",
					},
				},
			},
		},
		{
			name:  "valid integer literal",
			input: "5;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		{
			name:  "prefix bang",
			input: "!5;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.BANG, Literal: "!"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.BANG, Literal: "!"},
						Operator: "!",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "prefix minus",
			input: "-5;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.MINUS, Literal: "-"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.MINUS, Literal: "-"},
						Operator: "-",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix addition",
			input: "5 + 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.PLUS, Literal: "+"},
						Operator: "+",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix subtraction",
			input: "5 - 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.MINUS, Literal: "-"},
						Operator: "-",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix multiplication",
			input: "5 * 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.STAR, Literal: "*"},
						Operator: "*",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix division",
			input: "5 / 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.FSLASH, Literal: "/"},
						Operator: "/",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix equals",
			input: "5 == 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.EQ, Literal: "=="},
						Operator: "==",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix not equals",
			input: "5 != 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.NE, Literal: "!="},
						Operator: "!=",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "valid return statement",
			input: "return 5;",
			expected: []ast.Statement{
				&ast.ReturnStatement{
					Token: token.Token{Type: token.RETURN, Literal: "return"},
					Value: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		{
			name:  "valid let statement",
			input: "let x = 5;",
			expected: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		{
			name:  "invalid let statement: missing assignment operator",
			input: "let x 5;",
			expectedErrs: []string{
				"expected next token to be =, got INT instead",
			},
		},
		{
			name:  "invalid let statement: missing identifier",
			input: "let = 10;",
			expectedErrs: []string{
				"expected next token to be IDENT, got = instead",
				"no prefix parser found for =",
			},
		},
		{
			name:         "invalid let statement: missing identifier and assigner",
			input:        "let 10;",
			expectedErrs: []string{"expected next token to be IDENT, got INT instead"},
		},
		{
			name:  "valid boolean literal",
			input: "true;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.TRUE, Literal: "true"},
					Expression: &ast.Boolean{
						Token: token.Token{Type: token.TRUE, Literal: "true"},
						Value: true,
					},
				},
			},
		},
		{
			name:  "conditional expression no else",
			input: "if (x > 5) { x }",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{Type: token.IF, Literal: "if"},
						Condition: &ast.InfixExpression{
							Token: token.Token{Type: token.GT, Literal: ">"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Operator: ">",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
						Consequence: &ast.BlockStatement{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "conditional expression with else",
			input: "if (x > 5) { x } else { y }",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{Type: token.IF, Literal: "if"},
						Condition: &ast.InfixExpression{
							Token: token.Token{Type: token.GT, Literal: ">"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Operator: ">",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
						Consequence: &ast.BlockStatement{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
								},
							},
						},
						Alternative: &ast.BlockStatement{
							Token: token.Token{Type: token.IDENT, Literal: "y"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "y"},
									Expression: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "y"},
										Value: "y",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "function statement: multiple params",
			input: "fn add(x, y) { x + y; }",
			expected: []ast.Statement{
				&ast.FunctionStatement{
					Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
					Name:  ast.NewIdentifier("add"),
					Parameters: []*ast.Identifier{
						ast.NewIdentifier("x"),
						ast.NewIdentifier("y"),
					},
					Body: &ast.BlockStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Expression: &ast.InfixExpression{
									Token:    token.Token{Type: token.PLUS, Literal: "+"},
									Left:     ast.NewIdentifier("x"),
									Operator: "+",
									Right:    ast.NewIdentifier("y"),
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "function literal: one param",
			input: "fn echo(x) { x; }",
			expected: []ast.Statement{
				&ast.FunctionStatement{
					Token:      token.Token{Type: token.FUNCTION, Literal: "fn"},
					Name:       ast.NewIdentifier("echo"),
					Parameters: []*ast.Identifier{ast.NewIdentifier("x")},
					Body: &ast.BlockStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Token:      token.Token{Type: token.IDENT, Literal: "x"},
								Expression: ast.NewIdentifier("x"),
							},
						},
					},
				},
			},
		},
		{
			name:  "function literal: no params",
			input: "fn print5() { 5; }",
			expected: []ast.Statement{
				&ast.FunctionStatement{
					Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
					Name:  ast.NewIdentifier("print5"),
					Body: &ast.BlockStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Token:      token.Token{Type: token.INT, Literal: "5"},
								Expression: ast.NewIntegerLiteral(5),
							},
						},
					},
				},
			},
		},
		{
			name:  "function call expression: multiple args",
			input: "add(1, 2)",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "add"},
					Expression: &ast.CallExpression{
						Token: token.Token{Type: token.LPAREN, Literal: "("},
						Function: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "add"},
							Value: "add",
						},
						Arguments: []ast.Expression{
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2"},
								Value: 2,
							},
						},
					},
				},
			},
		},
		{
			name:  "function call expression: one arg",
			input: "add(1)",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "add"},
					Expression: &ast.CallExpression{
						Token: token.Token{Type: token.LPAREN, Literal: "("},
						Function: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "add"},
							Value: "add",
						},
						Arguments: []ast.Expression{
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
			},
		},
		{
			name:  "function call expression: no args",
			input: "add()",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "add"},
					Expression: &ast.CallExpression{
						Token: token.Token{Type: token.LPAREN, Literal: "("},
						Function: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "add"},
							Value: "add",
						},
					},
				},
			},
		},
		{
			name:  "strings",
			input: `"hello, world!"`,
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.STRING, Literal: "hello, world!"},
					Expression: &ast.StringLiteral{
						Token: token.Token{Type: token.STRING, Literal: "hello, world!"},
						Value: "hello, world!",
					},
				},
			},
		},
		{
			name:  "while statement: boolean condition",
			input: "while (true) { x + 1; }",
			expected: []ast.Statement{
				&ast.WhileStatement{
					Token: token.Token{Type: token.WHILE, Literal: "while"},
					Condition: &ast.Boolean{
						Token: token.Token{Type: token.TRUE, Literal: "true"},
						Value: true,
					},
					Block: &ast.BlockStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Expression: &ast.InfixExpression{
									Token: token.Token{Type: token.PLUS, Literal: "+"},
									Left: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
									Operator: "+",
									Right: &ast.IntegerLiteral{
										Token: token.Token{Type: token.INT, Literal: "1"},
										Value: 1,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "while statement: infix expression condition",
			input: "while (2 > 1) { x + 1; }",
			expected: []ast.Statement{
				&ast.WhileStatement{
					Token: token.Token{Type: token.WHILE, Literal: "while"},
					Condition: &ast.InfixExpression{
						Token: token.Token{Type: token.GT, Literal: ">"},
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "2"},
							Value: 2,
						},
						Operator: ">",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "1"},
							Value: 1,
						},
					},
					Block: &ast.BlockStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Expression: &ast.InfixExpression{
									Token: token.Token{Type: token.PLUS, Literal: "+"},
									Left: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
									Operator: "+",
									Right: &ast.IntegerLiteral{
										Token: token.Token{Type: token.INT, Literal: "1"},
										Value: 1,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "while statement: prefix expression condition",
			input: "while (!true) { x + 1; }",
			expected: []ast.Statement{
				&ast.WhileStatement{
					Token: token.Token{Type: token.WHILE, Literal: "while"},
					Condition: &ast.PrefixExpression{
						Token:    token.Token{Type: token.BANG, Literal: "!"},
						Operator: "!",
						Right: &ast.Boolean{
							Token: token.Token{Type: token.TRUE, Literal: "true"},
							Value: true,
						},
					},
					Block: &ast.BlockStatement{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Statements: []ast.Statement{
							&ast.ExpressionStatement{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Expression: &ast.InfixExpression{
									Token: token.Token{Type: token.PLUS, Literal: "+"},
									Left: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
									Operator: "+",
									Right: &ast.IntegerLiteral{
										Token: token.Token{Type: token.INT, Literal: "1"},
										Value: 1,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "assignment expression",
			input: "let x = 0; x = x + 1;",
			expected: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "0"},
						Value: 0,
					},
				},
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Expression: &ast.AssignmentExpression{
						Token: token.Token{Type: token.ASSIGN, Literal: "="},
						Left: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Right: &ast.InfixExpression{
							Token: token.Token{Type: token.PLUS, Literal: "+"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Operator: "+",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
			},
		},
		{
			name:  "list literal: same types",
			input: "[1, 2, 3]",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.LBRACK, Literal: "["},
					Expression: &ast.ListLiteral{
						Token: token.Token{Type: token.LBRACK, Literal: "["},
						Elems: []ast.Expression{
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2"},
								Value: 2,
							},
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "3"},
								Value: 3,
							},
						},
					},
				},
			},
		},
		{
			name:  "list literals: mixed types",
			input: `[1, "hello", fn(x) { x + 1 }]`,
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.LBRACK, Literal: "["},
					Expression: &ast.ListLiteral{
						Token: token.Token{Type: token.LBRACK, Literal: "["},
						Elems: []ast.Expression{
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
							&ast.StringLiteral{
								Token: token.Token{Type: token.STRING, Literal: "hello"},
								Value: "hello",
							},
							&ast.FunctionLiteral{
								Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
								Parameters: []*ast.Identifier{
									{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
								},
								Body: &ast.BlockStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Statements: []ast.Statement{
										&ast.ExpressionStatement{
											Token: token.Token{Type: token.IDENT, Literal: "x"},
											Expression: &ast.InfixExpression{
												Token: token.Token{Type: token.PLUS, Literal: "+"},
												Left: &ast.Identifier{
													Token: token.Token{
														Type:    token.IDENT,
														Literal: "x",
													},
													Value: "x",
												},
												Operator: "+",
												Right: &ast.IntegerLiteral{
													Token: token.Token{
														Type:    token.INT,
														Literal: "1",
													},
													Value: 1,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "index expression: typical",
			input: "[1][1]",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.LBRACK, Literal: "["},
					Expression: &ast.IndexExpression{
						Token: token.Token{Type: token.LBRACK, Literal: "["},
						Left: &ast.ListLiteral{
							Token: token.Token{Type: token.LBRACK, Literal: "["},
							Elems: []ast.Expression{
								&ast.IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "1"},
									Value: 1,
								},
							},
						},
						Index: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "1"},
							Value: 1,
						},
					},
				},
			},
		},
		{
			name:  "map expression: empty",
			input: `{}`,
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.LBRACE, Literal: "{"},
					Expression: &ast.MapLiteral{
						Token:   token.Token{Type: token.LBRACE, Literal: "{"},
						Entries: map[ast.Expression]ast.Expression{},
					},
				},
			},
		},
		{
			name:  "function expression: one parameter",
			input: "let add1 = fn(x) { x + 1 }",
			expected: []ast.Statement{
				&ast.LetStatement{
					Token: token.NewKeyword("let"),
					Name:  ast.NewIdentifier("add1"),
					Value: &ast.FunctionLiteral{
						Token: token.NewKeyword("fn"),
						Parameters: []*ast.Identifier{
							ast.NewIdentifier("x"),
						},
						Body: &ast.BlockStatement{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.InfixExpression{
										Token:    token.Token{Type: token.PLUS, Literal: "+"},
										Left:     ast.NewIdentifier("x"),
										Operator: "+",
										Right:    ast.NewIntegerLiteral(1),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "function expression: multi parameter",
			input: "let add = fn(x, y) { x + y }",
			expected: []ast.Statement{
				&ast.LetStatement{
					Token: token.NewKeyword("let"),
					Name:  ast.NewIdentifier("add"),
					Value: &ast.FunctionLiteral{
						Token: token.NewKeyword("fn"),
						Parameters: []*ast.Identifier{
							ast.NewIdentifier("x"),
							ast.NewIdentifier("y"),
						},
						Body: &ast.BlockStatement{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.InfixExpression{
										Token:    token.Token{Type: token.PLUS, Literal: "+"},
										Left:     ast.NewIdentifier("x"),
										Operator: "+",
										Right:    ast.NewIdentifier("y"),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "function statement: higher order",
			input: `
			fn adder(x) {
				fn applier(y) {
					x + y;
				}
			}
			`,
			expected: []ast.Statement{
				&ast.FunctionStatement{
					Token:      token.NewKeyword("fn"),
					Name:       ast.NewIdentifier("adder"),
					Parameters: []*ast.Identifier{ast.NewIdentifier("x")},
					Body: &ast.BlockStatement{
						Token: token.NewKeyword("fn"),
						Statements: []ast.Statement{
							&ast.FunctionStatement{
								Token:      token.NewKeyword("fn"),
								Name:       ast.NewIdentifier("applier"),
								Parameters: []*ast.Identifier{ast.NewIdentifier("y")},
								Body: &ast.BlockStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Statements: []ast.Statement{
										&ast.ExpressionStatement{
											Token: token.Token{Type: token.IDENT, Literal: "x"},
											Expression: &ast.InfixExpression{
												Token:    token.New(token.PLUS, '+'),
												Left:     ast.NewIdentifier("x"),
												Operator: "+",
												Right:    ast.NewIdentifier("y"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "class statement: empty body",
			input: "class Item {}",
			expected: []ast.Statement{
				&ast.ClassStatement{
					Token:   token.NewKeyword("class"),
					Name:    ast.NewIdentifier("Item"),
					Methods: nil,
				},
			},
		},
		{
			name: "class statement: one method",
			input: `
			class Item {
				add1(x) {
					return x + 1;
				}
			}`,
			expected: []ast.Statement{
				&ast.ClassStatement{
					Token: token.NewKeyword("class"),
					Name:  ast.NewIdentifier("Item"),
					Methods: []*ast.FunctionStatement{
						{
							Token:      token.NewKeyword("fn"),
							Name:       ast.NewIdentifier("add1"),
							Parameters: []*ast.Identifier{ast.NewIdentifier("x")},
							Body: &ast.BlockStatement{
								Token: token.NewKeyword("return"),
								Statements: []ast.Statement{
									&ast.ReturnStatement{
										Token: token.NewKeyword("return"),
										Value: &ast.InfixExpression{
											Token:    token.New(token.PLUS, '+'),
											Left:     ast.NewIdentifier("x"),
											Operator: "+",
											Right:    ast.NewIntegerLiteral(1),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "dot operator: <expression>.<identifier>",
			input: "object.property",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.NewIdent("object"),
					Expression: &ast.GetExpression{
						Token: token.NewSpecial("."),
						Left:  ast.NewIdentifier("object"),
						Right: ast.NewIdentifier("property"),
					},
				},
			},
		},
		{
			name:  "dot operator: <call>.<identifer>",
			input: "add().property",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.NewIdent("add"),
					Expression: &ast.GetExpression{
						Token: token.NewSpecial("."),
						Left: &ast.CallExpression{
							Token:    token.NewSpecial("("),
							Function: ast.NewIdentifier("add"),
						},
						Right: ast.NewIdentifier("property"),
					},
				},
			},
		},
		{
			name:  "dot operator: <identifier>.<call>",
			input: "object.add()",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.NewIdent("object"),
					Expression: &ast.CallExpression{
						Token: token.NewSpecial(token.LPAREN),
						Function: &ast.GetExpression{
							Token: token.NewSpecial(token.DOT),
							Left:  ast.NewIdentifier("object"),
							Right: ast.NewIdentifier("add"),
						},
					},
				},
			},
		},
		{
			name:  "dot operator: chaining",
			input: "parent.child.subchild",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.NewIdent("parent"),
					Expression: ast.NewGetExpression(
						ast.NewGetExpression(
							ast.NewIdentifier("parent"),
							ast.NewIdentifier("child"),
						),
						ast.NewIdentifier("subchild"),
					),
				},
			},
		},
		{
			name:  "dot operator: assignment",
			input: "parent.child = 1",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.NewIdent("parent"),
					Expression: &ast.AssignmentExpression{
						Token: token.NewSpecial(token.ASSIGN),
						Left: &ast.GetExpression{
							Token: token.NewSpecial(token.DOT),
							Left:  ast.NewIdentifier("parent"),
							Right: ast.NewIdentifier("child"),
						},
						Right: ast.NewIntegerLiteral(1),
					},
				},
			},
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			assert.NotNil(t, program)
			if tc.expectedErrs == nil {
				assert.Nil(t, p.Errors())
			} else {
				assert.Equal(t, tc.expectedErrs, p.Errors())
			}

			if tc.expected != nil {
				assert.EqualValues(t, tc.expected, program.Statements)
			}
		})
	}
}

func TestParser_MapLiteral(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected *ast.MapLiteral
	}{
		{
			name:  "map expression: literal",
			input: `{"a": "b", 1: "one", "b": true}`,
			expected: &ast.MapLiteral{
				Token: token.Token{Type: token.LBRACE, Literal: "{"},
				Entries: map[ast.Expression]ast.Expression{
					&ast.StringLiteral{
						Token: token.Token{Type: token.STRING, Literal: "a"},
						Value: "a",
					}: &ast.StringLiteral{
						Token: token.Token{Type: token.STRING, Literal: "b"},
						Value: "b",
					},
					&ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "1"},
						Value: 1,
					}: &ast.StringLiteral{
						Token: token.Token{Type: token.STRING, Literal: "one"},
						Value: "one",
					},
					&ast.StringLiteral{
						Token: token.Token{Type: token.STRING, Literal: "b"},
						Value: "b",
					}: &ast.Boolean{
						Token: token.Token{Type: token.TRUE, Literal: "true"},
						Value: true,
					},
				},
			},
		},
	}

	for _, testCase := range cases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			stmt := program.Statements[0]
			s, ok := stmt.(*ast.ExpressionStatement)
			assert.True(t, ok)
			m, ok := s.Expression.(*ast.MapLiteral)
			assert.True(t, ok)

			assert.Equal(t, tc.expected.Token, m.Token)
			assert.ElementsMatch(t, maps.Keys(tc.expected.Entries), maps.Keys(m.Entries))
			assert.ElementsMatch(t, maps.Values(tc.expected.Entries), maps.Values(m.Entries))
		})
	}
}

func TestParser_OperatorPrecedence(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unary mixed with binary",
			input:    "-a * b",
			expected: "((-a) * b)",
		},
		{
			name:     "multiple unary",
			input:    "!-a",
			expected: "(!(-a))",
		},
		{
			name:     "multiple addition",
			input:    "a + b + c",
			expected: "((a + b) + c)",
		},
		{
			name:     "addition and subtraction",
			input:    "a + b - c",
			expected: "((a + b) - c)",
		},
		{
			name:     "multiple multiplication",
			input:    "a * b * c",
			expected: "((a * b) * c)",
		},
		{
			name:     "multiplication and division",
			input:    "a * b / c",
			expected: "((a * b) / c)",
		},
		{
			name:     "multiple in sequence",
			input:    "a - b * c / d * e + f",
			expected: "((a - (((b * c) / d) * e)) + f)",
		},
		{
			name:     "multiple statements",
			input:    "a + b; x * y",
			expected: "(a + b)(x * y)",
		},
		{
			name:     "comparison operations",
			input:    "a > b == c > d",
			expected: "((a > b) == (c > d))",
		},
		{
			name:     "comparison with unary and binary operations",
			input:    "a * -b != -c / d",
			expected: "((a * (-b)) != ((-c) / d))",
		},
		{
			name:     "comparison with boolean literals",
			input:    "true != false",
			expected: "(true != false)",
		},
		{
			name:     "mixed numeric, identifiers, and boolean literals",
			input:    "3 > a != true",
			expected: "((3 > a) != true)",
		},
		{
			name:     "grouped expression",
			input:    "(5 + 5) * 2",
			expected: "((5 + 5) * 2)",
		},
		{
			name:     "call expression",
			input:    "add(1, 2 + 3, 4 * 5 + 6)",
			expected: "add(1, (2 + 3), ((4 * 5) + 6))",
		},
		{
			name:     "index expression",
			input:    "a + [1, 2, 3][4] + b",
			expected: "((a + ([1, 2, 3][4])) + b)",
		},
		{
			name:     "get expression",
			input:    "object.property.method()",
			expected: "((object.property).method)()",
		},
	}

	for _, testCase := range cases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expected, program.String())
		})
	}
}
