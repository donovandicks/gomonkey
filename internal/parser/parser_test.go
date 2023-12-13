package parser_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/parser"
	"github.com/donovandicks/gomonkey/internal/token"
	"github.com/stretchr/testify/assert"
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
			name:  "function literal: multiple params",
			input: "fn(x, y) { x + y; }",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
					Expression: &ast.FunctionLiteral{
						Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
						Parameters: []*ast.Identifier{
							{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
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
											Token: token.Token{Type: token.IDENT, Literal: "x"},
											Value: "x",
										},
										Operator: "+",
										Right: &ast.Identifier{
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
		},
		{
			name:  "function literal: one param",
			input: "fn(x) { x; }",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
					Expression: &ast.FunctionLiteral{
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
			name:  "function literal: no params",
			input: "fn() { 5; }",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
					Expression: &ast.FunctionLiteral{
						Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
						Body: &ast.BlockStatement{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Expression: &ast.IntegerLiteral{
										Token: token.Token{Type: token.INT, Literal: "5"},
										Value: 5,
									},
								},
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
